package model

import (
	"testing"

	"github.com/franela/goblin"
)

func TestGroupedTasks(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("GroupedTasks", func() {
		g.It("TaskCount() should return total count of tasks", func() {
			gt := GroupedTasks{
				Groups: []Group{
					Group{
						Name: "alpha",
						Tasks: []Task{
							Task{
								Spec: "alpha_spec_01.spec.ts",
								Viewport: "desktop",
								Browser: "chrome",
							},
							Task{
								Spec: "alpha_spec_01.spec.ts",
								Viewport: "mobile",
								Browser: "chrome",
							},
						},
					},
					Group{
						Name: "beta",
						Tasks: []Task{
							Task{
								Spec: "beta_spec_01.spec.ts",
								Viewport: "mobile",
								Browser: "webkit",
							},
						},
					},
				},
			}

			g.Assert(gt.TaskCount()).Eql(3)
		})
	})
}