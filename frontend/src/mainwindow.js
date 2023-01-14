import React, {useState, useEffect} from 'react'

import io from 'socket.io-client';

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import CodeWindow from './codewindow';

const CREATENEWROOMAPI = 'http://localhost:8080/api/v1/room/create'
const BACKENDSERVER = 'http://localhost:8080/'

const socket = io("http://127.0.0.1:8080/");
// client-side
socket.on("connect", () => {
  console.log(socket.id);
});

function MainWindow() {
  const [isLogin, setLogin] = useState(false);
  const [room, setRoom] = useState(null);
  const [name, setName] = useState("");

  const CreateNewRoom = async () => {
    console.log("Creating room...")
    const response = await fetch(CREATENEWROOMAPI, {
      method: 'POST'
    })
    const data = await response.json()
    console.log(data)
    setRoom(data.room_id)
    setLogin(true)
    socket.emit("joinroom", data.room_id)
  }

  const JoinRoom = async () => {
    console.log("joining room")
    socket.emit("joinroom", room)
  }

  useEffect(() => {

  }, [isLogin])

  return (
    <div style={{background: '#C5D4D2', height: '100vh'}}>
      { 
        isLogin ? <CodeWindow room_id={room}/> :
        (
          <div className='Container'>
            <div className='row'>
              <h1 className='my-5'>CodeBooks</h1>
            </div>
            <div className='row'>
              <div className='col' />
              <div className='col' >
                <div className='p-4 mt-5' style={{background: '#17C3A5', borderRadius: 10}}>
                  <Form.Control
                    id="pin"
                    className='my-2'
                    placeholder='Room PIN'
                    onChange={(e) => setRoom(e.target.value)}
                  />
                  <Form.Control
                    id="name"
                    className='my-2'
                    placeholder='Nickname'
                    onChange={(e) => setName(e.target.value)}
                  />
                  <div className='row mb-2 mt-5'>
                    <div className='col-2' />
                    <div className='col'>
                      <Button onClick={() => JoinRoom()}>Enter</Button>
                    </div>
                    <div className='col-2' />
                  </div>
                  <div className='row'>
                    <div className='col-2' />
                      <div className='col'>
                        <Button onClick={() => CreateNewRoom()}>Create New Room</Button>
                      </div>
                    <div className='col-2' />
                  </div>
                </div>
              </div>
              <div className='col' />
            </div>
          </div>
        )
      }
    </div>
  )
  
}

export default MainWindow