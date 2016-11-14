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

package xplane

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ornen/go-xplane/messages"
	"log"
	"net"
)

const (
	datagramPrefixLength = 5
	messageLength        = 36
)

type XPlane struct {
	RemoteAddress string
	LocalAddress  string
	Messages      chan Message
	connection    *net.UDPConn
}

func New(remoteAddress, localAddress string) XPlane {
	return XPlane{
		RemoteAddress: remoteAddress,
		LocalAddress:  localAddress,
		Messages:      make(chan Message),
	}
}

func (x *XPlane) Receive() {
	serverAddr, err := net.ResolveUDPAddr("udp", x.LocalAddress)
	serverConn, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		panic(err)
	}

	defer serverConn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, _ := serverConn.ReadFromUDP(buf)
		m := (n - datagramPrefixLength) / messageLength

		for i := 0; i < m; i++ {
			sentence := buf[datagramPrefixLength+i*messageLength : datagramPrefixLength+(i+1)*messageLength]
			x.parse(sentence)
		}
	}
}

func (x *XPlane) parse(sentence []byte) {
	messageType := sentence[0]
	messageBuffer := bytes.NewBuffer(sentence[4:])

	messageData := make([]float32, 8)
	binary.Read(messageBuffer, binary.LittleEndian, &messageData)

	switch messageType {
	case messages.SpeedMessageType:
		x.Messages <- messages.NewSpeedMessage(messageData)
	case messages.GLoadMessageType:
		x.Messages <- messages.NewGLoadMessage(messageData)
	case messages.AngularVelocitiesMessageType:
		x.Messages <- messages.NewAngularVelocitiesMessage(messageData)
	case messages.PitchRollHeadingMessageType:
		x.Messages <- messages.NewPitchRollHeadingMessage(messageData)
	case messages.FlightControlMessageType:
		x.Messages <- messages.NewFlightControlMessage(messageData)
	case messages.GearsBrakesMessageType:
		x.Messages <- messages.NewGearsBrakesMessage(messageData)
	case messages.WeatherMessageType:
		x.Messages <- messages.NewWeatherMessage(messageData)
	case messages.LatLonAltMessageType:
		x.Messages <- messages.NewLatLonAltMessage(messageData)
	case messages.LocVelDistTraveledMessageType:
		x.Messages <- messages.NewLocVelDistTraveledMessage(messageData)
	case messages.BatteryAmperageMessageType:
		x.Messages <- messages.NewBatteryAmperageMessage(messageData)
	case messages.BatteryVoltageMessageType:
		x.Messages <- messages.NewBatteryVoltageMessage(messageData)
	case messages.EngineRPMMessageType:
		x.Messages <- messages.NewEngineRPMMessage(messageData)
	case messages.PropRPMMessageType:
		x.Messages <- messages.NewPropRPMMessage(messageData)
	case messages.PropPitchMessageType:
		x.Messages <- messages.NewPropPitchMessage(messageData)
	default:
		log.Println("Unknown message type: ", messageType)
	}
}

func (x *XPlane) Connect() {
	udpAddr, err := net.ResolveUDPAddr("udp", x.RemoteAddress)

	if err != nil {
		fmt.Println("Wrong address!")
		return
	}

	x.connection, err = net.DialUDP("udp", nil, udpAddr)
}

func (x *XPlane) Send(command Command) {
	commandData := command.Data()

	buf := new(bytes.Buffer)

	buf.Write([]byte{'D', 'A', 'T', 'A', 0})
	buf.Write([]byte{byte(command.Type()), 0, 0, 0})

	if err := binary.Write(buf, binary.LittleEndian, &commandData); err != nil {
		fmt.Println(err)
		return
	}

	x.connection.Write(buf.Bytes())
}
