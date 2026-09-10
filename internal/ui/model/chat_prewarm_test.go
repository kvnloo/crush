package model

import (
	"testing"

	"github.com/charmbracelet/crush/internal/ui/chat"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/list"
	"github.com/stretchr/testify/require"
)

// mockMessageItem is a test helper that implements chat.MessageItem.
type mockMessageItem struct {
	*list.Versioned
	id      string
	content string
}

func newMockMessageItem(id, content string) *mockMessageItem {
	return &mockMessageItem{
		Versioned: list.NewVersioned(),
		id:        id,
		content:   content,
	}
}

func (m *mockMessageItem) ID() string { return m.id }

func (m *mockMessageItem) Render(width int) string {
	return m.content
}

func (m *mockMessageItem) RawRender(width int) string {
	return m.content
}

func (m *mockMessageItem) Finished() bool { return true }

// TestChat_SetMessages_StartsPrewarm covers the PER-877 fix: SetMessages
// must start the prewarm sequence so that the first scrollbar draw (which
// needs Offset or TotalHeight) doesn't walk the entire list. Without this,
// loading a large session causes a ~30s hang on first scroll/render.
func TestChat_SetMessages_StartsPrewarm(t *testing.T) {
	t.Parallel()

	// Create a chat with a large session (simulate 100 messages).
	com := &common.Common{}
	c := NewChat(com, "default")
	c.SetSize(80, 24)

	msgs := make([]chat.MessageItem, 100)
	for i := range msgs {
		msgs[i] = newMockMessageItem("msg-"+string(rune('0'+i)), "content")
	}

	// Call SetMessages and capture the returned command.
	cmd := c.SetMessages(msgs...)
	require.NotNil(t, cmd, "SetMessages must return a prewarm command")

	// The command should be a chatWarmMsg with seq 1 (first warm cycle).
	// Execute it and verify it starts the warming process.
	msg := cmd()
	warmMsg, ok := msg.(chatWarmMsg)
	require.True(t, ok, "SetMessages command must be a chatWarmMsg")
	require.Equal(t, 1, warmMsg.seq, "first warm sequence should be 1")

	// Verify that resizing flag is set.
	require.True(t, c.resizing, "SetMessages must set resizing=true to suppress scrollbar")
	require.Equal(t, 0, c.warmNext, "warmNext must start at 0")
	require.Equal(t, 1, c.resizeSettleSeq, "resizeSettleSeq must be 1 after first SetMessages")

	// Execute one warming step to verify the mechanism works.
	cmd, done := c.WarmStep(warmMsg.seq)
	require.NotNil(t, cmd, "first warm step must return a continuation command")
	require.False(t, done, "first warm step must not be done")
	require.Equal(t, warmBatchSize, c.warmNext, "warmNext must advance by warmBatchSize")
}

// TestChat_SetMessages_PrewarmCompletesEventually verifies that the
// prewarm sequence runs to completion and clears the resizing flag.
func TestChat_SetMessages_PrewarmCompletesEventually(t *testing.T) {
	t.Parallel()

	com := &common.Common{}
	c := NewChat(com, "default")
	c.SetSize(80, 24)

	// Create a small session that can be warmed in a few steps.
	msgs := make([]chat.MessageItem, warmBatchSize*2+5)
	for i := range msgs {
		msgs[i] = newMockMessageItem("msg-"+string(rune('0'+i)), "content")
	}

	cmd := c.SetMessages(msgs...)
	require.NotNil(t, cmd)

	msg := cmd()
	warmMsg, ok := msg.(chatWarmMsg)
	require.True(t, ok)

	// Run warming steps until completion.
	var step int
	for {
		cmd, done := c.WarmStep(warmMsg.seq)
		step++
		if done {
			require.Nil(t, cmd, "final warm step must return nil command")
			require.False(t, c.resizing, "resizing must be cleared on completion")
			require.Equal(t, len(msgs), c.warmNext, "warmNext must reach item count")
			break
		}
		require.NotNil(t, cmd, "intermediate warm step must return a command")
		msg = cmd()
		warmMsg, ok = msg.(chatWarmMsg)
		require.True(t, ok)
		require.Less(t, step, 10, "warming should complete in a reasonable number of steps")
	}
}

// TestChat_SetMessages_SupportsLargeSessions verifies that SetMessages
// handles sessions with thousands of messages without blocking.
func TestChat_SetMessages_SupportsLargeSessions(t *testing.T) {
	t.Parallel()

	com := &common.Common{}
	c := NewChat(com, "default")
	c.SetSize(80, 24)

	// Create a large session (5000 messages, similar to the issue report).
	const messageCount = 5000
	msgs := make([]chat.MessageItem, messageCount)
	for i := range msgs {
		msgs[i] = newMockMessageItem("msg-"+string(rune('0'+i%100)), "content line\nmore content")
	}

	cmd := c.SetMessages(msgs...)
	require.NotNil(t, cmd, "SetMessages must return a prewarm command even for large sessions")

	// Verify the warming starts correctly.
	msg := cmd()
	warmMsg, ok := msg.(chatWarmMsg)
	require.True(t, ok)
	require.Equal(t, 1, warmMsg.seq)

	// Run a few warming steps to verify the incremental warming works.
	for range 5 {
		cmd, done := c.WarmStep(warmMsg.seq)
		if done {
			break
		}
		require.NotNil(t, cmd)
		msg = cmd()
		warmMsg, ok = msg.(chatWarmMsg)
		require.True(t, ok)
	}
	// Warming should not complete in just 5 steps for 5000 messages (warmBatchSize=25).
	require.True(t, c.resizing, "resizing should still be active after a few steps")
	require.Less(t, c.warmNext, messageCount, "warming should not be complete yet")
}
