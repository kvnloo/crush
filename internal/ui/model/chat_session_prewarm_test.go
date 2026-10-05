package model

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/crush/internal/ui/chat"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/list"
	"github.com/stretchr/testify/require"
)

type prewarmMessageItem struct {
	*list.Versioned
	id      string
	content string
}

func newPrewarmMessageItem(id, content string) *prewarmMessageItem {
	return &prewarmMessageItem{Versioned: list.NewVersioned(), id: id, content: content}
}

func (m *prewarmMessageItem) ID() string           { return m.id }
func (m *prewarmMessageItem) Render(int) string    { return m.content }
func (m *prewarmMessageItem) RawRender(int) string { return m.content }
func (m *prewarmMessageItem) Finished() bool       { return true }

func TestChatSetMessagesStartsIncrementalPrewarm(t *testing.T) {
	t.Parallel()

	c := NewChat(&common.Common{}, "default")
	c.SetSize(80, 24)

	msgs := make([]chat.MessageItem, warmBatchSize*2+5)
	for i := range msgs {
		msgs[i] = newPrewarmMessageItem(fmt.Sprintf("msg-%d", i), "content")
	}

	cmd := c.SetMessages(msgs...)
	require.NotNil(t, cmd)
	warm, ok := cmd().(chatWarmMsg)
	require.True(t, ok)
	require.True(t, c.resizing)
	require.Zero(t, c.warmNext)

	cmd, done := c.WarmStep(warm.seq)
	require.False(t, done)
	require.NotNil(t, cmd)
	require.Equal(t, warmBatchSize, c.warmNext)

	for !done {
		msg := cmd()
		warm, ok = msg.(chatWarmMsg)
		require.True(t, ok)
		cmd, done = c.WarmStep(warm.seq)
	}
	require.False(t, c.resizing)
	require.Equal(t, len(msgs), c.warmNext)
}

func TestChatSetMessagesLargeSessionDoesNotWarmSynchronously(t *testing.T) {
	t.Parallel()

	c := NewChat(&common.Common{}, "default")
	c.SetSize(80, 24)

	const messageCount = 5000
	msgs := make([]chat.MessageItem, messageCount)
	for i := range msgs {
		msgs[i] = newPrewarmMessageItem(fmt.Sprintf("msg-%d", i), "content line\nmore content")
	}

	cmd := c.SetMessages(msgs...)
	require.NotNil(t, cmd)
	require.True(t, c.resizing)
	require.Zero(t, c.warmNext, "SetMessages should schedule warming, not synchronously render all items")
}
