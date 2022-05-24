package groupedtaskloader

import (
	"testing"

	"github.com/bspain/funkytown/shared/model"
	"github.com/franela/goblin"
)

func TestLoadGroupedTaskFile(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("LoadGroupedTaskFile", func() {
		g.It("should load a groupedtasklist from a file", func() {
			file := "../../../__fixtures__/group_tasks_file.json"

			groupedtasks, err := LoadGroupTasksFile(file)

			g.Assert(err == nil).IsTrue()
			g.Assert(len(groupedtasks.Groups)).Eql(2)

			{
				gp := groupedtasks.Groups[0]
				g.Assert(gp.Name).Eql("alpha")

				// Expect alpha_spec_01 to have 7 tasks (as mobile-firefox is an invalid combo)
				g.Assert((len(gp.Tasks))).Eql(7)		
				g.Assert(gp.Tasks[0]).Eql(model.Task{
					Spec: "alpha_spec_01.spec.ts",
					Viewport: "desktop",
					Browser: "chrome",
				})		
				g.Assert(gp.Tasks[1]).Eql(model.Task{
					Spec: "alpha_spec_01.spec.ts",
					Viewport: "desktop",
					Browser: "firefox",
				})		
				g.Assert(gp.Tasks[2]).Eql(model.Task{
					Spec: "alpha_spec_01.spec.ts",
					Viewport: "desktop",
					Browser: "webkit",
				})		
				g.Assert(gp.Tasks[3]).Eql(model.Task{
					Spec: "alpha_spec_01.spec.ts",
					Viewport: "mobile",
					Browser: "chrome",
				})		
				g.Assert(gp.Tasks[4]).Eql(model.Task{
					Spec: "alpha_spec_01.spec.ts",
					Viewport: "mobile",
					Browser: "webkit",
				})		
				g.Assert(gp.Tasks[5]).Eql(model.Task{
					Spec: "alpha_spec_02.spec.ts",
					Viewport: "mobile",
					Browser: "chrome",
				})		
				g.Assert(gp.Tasks[6]).Eql(model.Task{
					Spec: "alpha_spec_02.spec.ts",
					Viewport: "mobile",
					Browser: "webkit",
				})		
			}
			{
				gp := groupedtasks.Groups[1]
				g.Assert(gp.Name).Eql("beta")
				g.Assert((len(gp.Tasks))).Eql(1)		

				g.Assert(gp.Tasks[0]).Eql(model.Task{
					Spec: "beta_spec_01.spec.ts",
					Viewport: "mobile",
					Browser: "webkit",
				})		
			}
		})
	})
}