package daemon

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/env"
)

type TodoLine struct {
	Action string
	Commit *models.Commit
}

func (self *TodoLine) ToString() string {
	if self.Action == "break" {
		return self.Action + "\n"
	} else {
		return self.Action + " " + self.Commit.Sha + " " + self.Commit.Name + "\n"
	}
}

func TodoLinesToString(todoLines []TodoLine) string {
	lines := slices.Map(todoLines, func(todoLine TodoLine) string {
		return todoLine.ToString()
	})

	return strings.Join(slices.Reverse(lines), "")
}

type ChangeTodoAction struct {
	Sha       string
	NewAction todo.TodoCommand
}

func handleInteractiveRebase(common *common.Common, f func(path string) error) error {
	common.Log.Info("Lazygit invoked as interactive rebase demon")
	common.Log.Info("args: ", os.Args)
	path := os.Args[1]

	if strings.HasSuffix(path, "git-rebase-todo") {
		return f(path)
	} else if strings.HasSuffix(path, filepath.Join(gitDir(), "COMMIT_EDITMSG")) { // TODO: test
		// if we are rebasing and squashing, we'll see a COMMIT_EDITMSG
		// but in this case we don't need to edit it, so we'll just return
	} else {
		common.Log.Info("Lazygit demon did not match on any use cases")
	}

	return nil
}

func gitDir() string {
	dir := env.GetGitDirEnv()
	if dir == "" {
		return ".git"
	}
	return dir
}
