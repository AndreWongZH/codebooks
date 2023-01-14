import React, {useState, useEffect} from 'react'

import io from 'socket.io-client';

import Button from 'react-bootstrap/Button';
import CodeWindow from './codewindow';

const CREATENEWROOMAPI = 'http://localhost:8080/api/v1/room/create'
const BACKENDSERVER = 'http://localhost:8080'

const CreateNewRoom = async (setLogin) => {
  console.log("Creating room...")
  // const response = await fetch(CREATENEWROOMAPI, {
  //   method: 'POST'
  // })
  // const data = await response.json()
  // console.log(data)
  setLogin(true)
  

  // const socket = io(BACKENDSERVER);
}

function MainWindow() {
  const [isLogin, setLogin] = useState(false);

  useEffect(() => {

  }, [isLogin])

  return (
    <div>
      { 
        isLogin ? <CodeWindow /> :
        (
          <div className='Container'>
            <div className='row'>
              <div className='col'>
                <h1>CodeBooks</h1>
                <Button onClick={() => CreateNewRoom(setLogin)}>Create New Room</Button>
              </div>
            </div>
          </div>
        )
      }
    </div>
  )
  
}

export default MainWindow