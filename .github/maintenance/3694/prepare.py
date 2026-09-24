#!/usr/bin/env python3
"""One-shot, pinned integration for PR #3694; never updates a remote ref."""
import json
import os
from pathlib import Path
import subprocess
import sys

OLD = 'fda19ac199ed405c6a57203b571585c9002912e7'
BASE = '959ebd5d80805c3ec98d3ea49d1313e401b4e300'
MERGE_BASE = '559ec80922fecf3baa0b7599230f4c91067440de'
BRANCH = 'fix/3671-headless-unresolved-large'
PATHS = [
    'internal/app/app.go', 'internal/cmd/run.go', 'internal/cmd/run_model_test.go',
    'internal/config/config.go', 'internal/config/load.go',
    'internal/config/load_test.go', 'internal/config/store.go',
]
TEST = 'internal/cmd/run_refresh_test.go'
OLD_REFRESH = '\t\tif err != nil {\n\t\t\tslog.Debug("failed to refresh config after model override", "error", err)\n\t\t} else {\n\t\t\tws.Config = cfg\n\t\t}\n'
NEW_REFRESH = '\t\tif err != nil {\n\t\t\treturn fmt.Errorf("failed to refresh config after model override: %w", err)\n\t\t}\n\t\tws.Config = cfg\n'


def git(*args, input=None):
    return subprocess.check_output(['git', *args], input=input, text=True)


def replace_once(text, old, new, label):
    count = text.count(old)
    if count != 1:
        raise RuntimeError(f'{label}: expected one source anchor, found {count}')
    return text.replace(old, new, 1)


def changed(base):
    return set(git('diff', '--name-only', base).splitlines())


def assemble(assets, evidence):
    remote = git('ls-remote', 'origin', f'refs/heads/{BRANCH}').split()[0]
    if remote != OLD:
        raise RuntimeError(f'PR moved: expected {OLD}, found {remote}; refusing to overwrite')
    git('fetch', 'origin', f'refs/heads/{BRANCH}')
    git('fetch', 'https://github.com/charmbracelet/crush.git', BASE)
    if git('merge-base', OLD, BASE).strip() != MERGE_BASE:
        raise RuntimeError('Unexpected merge base')
    original_paths = set(git('diff', '--name-only', MERGE_BASE, OLD).splitlines())
    if original_paths != set(PATHS):
        raise RuntimeError(f'Unexpected original PR paths: {original_paths}')
    patch = git('diff', '--binary', MERGE_BASE, OLD)
    (evidence / 'original.patch').write_text(patch)
    old_run = git('show', f'{OLD}:internal/cmd/run.go')
    git('checkout', '--detach', BASE)
    # Preserve upstream plus the existing PR hunks in these five files.
    # Reconcile the two run paths explicitly with the new reasoning-effort
    # block; never use a blanket ours/theirs conflict resolution.
    git('apply', '--exclude=internal/cmd/run.go', '--exclude=internal/app/app.go', '-', input=patch)
    path = Path('internal/cmd/run.go')
    source = path.read_text()
    for cfg, indent in [('ws.Config', '\t\t\t'), ('ws.Config()', '\t\t')]:
        anchor = f'{indent}if !{cfg}.IsConfigured() {{\n{indent}\treturn fmt.Errorf("no providers configured - please run \'crush\' to set up a provider interactively")\n{indent}}}\n'
        addition = f'\n{indent}if err := refuseUnresolvedLarge(largeModel, {cfg}); err != nil {{\n{indent}\treturn err\n{indent}}}\n'
        source = replace_once(source, anchor, anchor + addition, f'{cfg} headless guard')
    override = '\tif largeModel != "" || smallModel != "" {\n\t\tif err := overrideModels(ctx, c, ws, largeModel, smallModel); err != nil {\n\t\t\treturn fmt.Errorf("failed to override models: %w", err)\n\t\t}\n\t}\n'
    baseline_override = override[:-3] + '\t\tcfg, err := c.GetConfig(ctx, ws.ID)\n' + OLD_REFRESH + '\t}\n'
    source = replace_once(source, override, baseline_override, 'model refresh')
    spinner = '\tvar (\n\t\tspinner   *format.Spinner\n'
    source = replace_once(source, spinner, '\tfmt.Fprintln(os.Stderr, resolvedLargeLine(ws.Config))\n\n' + spinner, 'client model pin')
    start = old_run.index('// refuseUnresolvedLarge fails crush run')
    end = old_run.index('// overrideModels resolves model strings', start)
    helpers = old_run[start:end]
    source = replace_once(source, '// overrideModels resolves model strings', helpers + '// overrideModels resolves model strings', 'existing shared helpers')
    path.write_text(source)
    path = Path('internal/app/app.go')
    path.write_text(replace_once(path.read_text(), spinner, '\tfmt.Fprintln(os.Stderr, app.config.Config().ResolvedLargeLine())\n\n' + spinner, 'in-process model pin'))
    subprocess.run(['gofmt', '-w', *PATHS], check=True)
    git('diff', '--check')
    if changed(BASE) != set(PATHS):
        raise RuntimeError('Integration changed unexpected files')
    git('add', '--', *PATHS)
    tree = git('write-tree').strip()
    merge = git('commit-tree', tree, '-p', OLD, '-p', BASE, '-m', 'chore: merge main into headless model validation').strip()
    git('reset', '--hard', merge)
    Path(TEST).write_text((assets / 'run_refresh_test.go.txt').read_text())
    subprocess.run(['gofmt', '-w', TEST], check=True)
    (evidence / 'integration.json').write_text(json.dumps({'old_head': OLD, 'upstream': BASE, 'merge': merge, 'tree': tree}, indent=2))


def fix():
    path = Path('internal/cmd/run.go')
    path.write_text(replace_once(path.read_text(), OLD_REFRESH, NEW_REFRESH, 'fail-loud refresh'))
    subprocess.run(['gofmt', '-w', str(path), TEST], check=True)
    git('diff', '--check')
    if changed('HEAD') != {'internal/cmd/run.go'}:
        # The new test is untracked until the final explicit add.
        raise RuntimeError('Unexpected tracked follow-up changes')


def verify_red(path):
    events = [json.loads(line) for line in Path(path).read_text().splitlines() if line.startswith('{')]
    outcomes = {e['Test']: e['Action'] for e in events if e.get('Test') and e['Action'] in ('pass', 'fail')}
    root = 'TestRunNonInteractive_ModelRefresh'
    expected = {f'{root}/http_error': 'fail', f'{root}/invalid_json': 'fail', f'{root}/success': 'pass'}
    if any(outcomes.get(name) != action for name, action in expected.items()):
        raise RuntimeError(f'Baseline did not demonstrate the expected bug and passing control: {outcomes}')
    output = ''.join(e.get('Output', '') for e in events)
    if 'failed to refresh config after model override' not in output or 'panic:' in output:
        raise RuntimeError('Baseline failed for an unrelated reason')


def finish(evidence):
    git('add', '--', 'internal/cmd/run.go', TEST)
    git('diff', '--cached', '--check')
    git('commit', '-m', 'fix(run): stop when post-override config refresh fails')
    if changed(BASE) != set(PATHS + [TEST]):
        raise RuntimeError('Candidate has unexpected changes relative to upstream')
    git('merge-base', '--is-ancestor', OLD, 'HEAD')
    git('merge-base', '--is-ancestor', BASE, 'HEAD')
    (evidence / 'candidate.patch').write_text(git('diff', BASE, 'HEAD'))
    (evidence / 'candidate.json').write_text(json.dumps({'head': git('rev-parse', 'HEAD').strip(), 'old_head': OLD, 'upstream': BASE}, indent=2))


if __name__ == '__main__':
    mode = sys.argv[1]
    assets = Path(os.environ['ASSETS'])
    evidence = Path(os.environ['EVIDENCE'])
    evidence.mkdir(parents=True, exist_ok=True)
    if mode == 'assemble':
        assemble(assets, evidence)
    elif mode == 'fix':
        fix()
    elif mode == 'verify-red':
        verify_red(evidence / 'baseline.jsonl')
    elif mode == 'finish':
        finish(evidence)
    else:
        raise SystemExit(f'Unknown mode: {mode}')
