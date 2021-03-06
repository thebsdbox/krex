// Copyright © 2018 Kris Nova <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"fmt"

	"strings"

	"github.com/manifoldco/promptui"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type List func(above *Vertex) (*Vertex, error)

func ListNamespaces(above *Vertex) (*Vertex, error) {
	ns, err := runtimeInstance.clientset.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	v := &Vertex{}
	prompt := promptui.Select{
		Label: "Select Namespace",
		Items: []string{""},
	}
	var list []string
	for _, n := range ns.Items {
		list = append(list, n.Name)
	}
	prompt.Items = list
	prompt.Size = len(list)
	v.Prompt = prompt
	v.Above = above
	v.ListFunc = ListApps
	return v, nil
}

func ListApps(above *Vertex) (*Vertex, error) {
	ss, err := runtimeInstance.clientset.AppsV1().StatefulSets(above.PreviousOutput).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	ds, err := runtimeInstance.clientset.AppsV1().Deployments(above.PreviousOutput).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	dss, err := runtimeInstance.clientset.AppsV1().DaemonSets(above.PreviousOutput).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	v := &Vertex{}
	prompt := promptui.Select{
		Label: "Select Applications",
		Items: []string{""},
	}
	var list []string
	for _, n := range ss.Items {
		list = append(list, fmt.Sprintf("[StatefulSet] %s", n.Name))
	}
	for _, n := range ds.Items {
		list = append(list, fmt.Sprintf("[Deployment] %s", n.Name))
	}
	for _, n := range dss.Items {
		list = append(list, fmt.Sprintf("[DaemonSets] %s", n.Name))
	}
	prompt.Items = list
	prompt.Size = len(list)
	v.Prompt = prompt
	v.Above = above
	v.ListFunc = ListPods
	v.Namespace = above.PreviousOutput
	return v, nil
}

func ListPods(above *Vertex) (*Vertex, error) {
	spl := strings.Split(above.PreviousOutput, " ")
	//kind := spl[0]
	name := spl[1]

	pods, err := runtimeInstance.clientset.CoreV1().Pods(above.Namespace).List(v1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"k8s-app": name,
		}).String(),
	})
	if err != nil {
		return nil, err
	}
	v := &Vertex{}
	prompt := promptui.Select{
		Label: "Select Pods",
		Items: []string{""},
	}

	var list []string
	for _, n := range pods.Items {
		//if strings.Contains(n.Name, name) {
		list = append(list, n.Name)
		//}
	}
	prompt.Items = list
	prompt.Size = len(list)
	v.Prompt = prompt
	v.Above = above
	v.ListFunc = ListPodActions
	v.Namespace = above.Namespace
	return v, nil
}

func ListPodActions(above *Vertex) (*Vertex, error) {
	v := &Vertex{}
	prompt := promptui.Select{
		Label: "Select Pod Action",
		Items: []string{"Edit", "Describe", "Logs", "Shell Debug"},
		Size:  4,
	}
	v.Prompt = prompt
	v.Above = above
	v.Namespace = above.Namespace
	v.ListFunc = ListEnd
	v.Terminate = true
	return v, nil
}

func ListEnd(above *Vertex) (*Vertex, error) {
	return nil, nil
}
