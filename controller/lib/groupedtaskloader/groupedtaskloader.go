package groupedtaskloader

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bspain/funkytown/shared/model"
)

type GroupTasksFile struct {
	Groups []GroupEntry	`json:"groups"`
}

type GroupEntry struct {
	Name string `json:"name"`
	Tasks []TaskMatrix `json:"tasks"`
}

type TaskMatrix struct {
	Spec	string `json:"spec"`
	Viewports []string `json:"viewports"`
	Browsers []string `json:"browsers"`
}

func LoadGroupTasksFile(file string) (model.GroupedTasks, error) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return model.GroupedTasks{}, err
	}

	var gtf GroupTasksFile
	err = json.Unmarshal(data, &gtf)

	// Convert file data with task matrix into flat task list
	gt := model.GroupedTasks{}
	for _, g := range gtf.Groups {
		group := model.Group{
			Name: g.Name,
		}

		for _, t := range g.Tasks {
			for _, v := range t.Viewports {
				for _, b := range t.Browsers {
					// browser: firefox , viewport: mobile is invalid
					if v == "mobile" && b == "firefox" {
						log.Printf("Skipping %s : %s : %s , as an invalid task matrix combination", t.Spec, b, v)
						continue
					}

					task := model.Task{
						Spec: t.Spec,
						Viewport: v,
						Browser: b,
					}
					group.Tasks = append(group.Tasks, task)
				}
			}
		}

		gt.Groups = append(gt.Groups, group)
	}

	return gt, nil
}