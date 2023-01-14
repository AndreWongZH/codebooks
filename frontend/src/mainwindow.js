import React, {useState, useEffect} from 'react'

import io from 'socket.io-client';

import Button from 'react-bootstrap/Button';
import CodeWindow from './codewindow';

const CREATENEWROOMAPI = 'http://localhost:8080/api/v1/room/create'
const BACKENDSERVER = 'http://localhost:8080/'

const socket = io("http://127.0.0.1:8080/");
// client-side
socket.on("connect", () => {
  console.log(socket.id); // x8WIv7-mJelg7on_ALbx
});

function MainWindow() {
  const [isLogin, setLogin] = useState(false);
  const [room, setRoom] = useState(null);

  const CreateNewRoom = async () => {
    console.log("Creating room...")
    const response = await fetch(CREATENEWROOMAPI, {
      method: 'POST'
    })
    const data = await response.json()
    console.log(data)
    setRoom(data.room_id)
    setLogin(true)
  }

  useEffect(() => {

  }, [isLogin])

  return (
    <div>
      { 
        isLogin ? <CodeWindow room_id={room}/> :
        (
          <div className='Container'>
            <div className='row'>
              <div className='col'>
                <h1>CodeBooks</h1>
                <Button onClick={() => CreateNewRoom(setLogin)}>Create New Room</Button>
                <Button onClick={() => socket.emit("result", "hello world")}>Emit</Button>
              </div>
            </div>
          </div>
        )
      }
    </div>
  )
  
}

export default MainWindow