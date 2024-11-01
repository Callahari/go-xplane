// Copyright 2016 Ornen. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package messages

type EngineThrust struct {
	MessageData
	Engine1 float64 //l,lb
	Engine2 float64
	Engine3 float64
	Engine4 float64
}

func NewEngineThrust(data []float32) EngineThrust {
	return EngineThrust{
		MessageData: MessageData{Name: "EngineThrust"},
		Engine1:     float64(data[0]),
		Engine2:     float64(data[1]),
		Engine3:     float64(data[2]),
		Engine4:     float64(data[3]),
	}
}

func (m EngineThrust) Type() uint { return EngineThrustType }

func (m EngineThrust) GetName() string { return m.Name }