import React, {useState, useEffect} from 'react'

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import CodeWindow from './codewindow';

const CREATENEWROOMAPI = 'http://localhost:8080/api/v1/room/create'
const CHECKROOMAPI = 'http://localhost:8080/api/v1/room/check?room_id='
const GETROOMAPI = 'http://localhost:8080/api/v1/room/get?room_id='

function MainWindow() {
  const [isLogin, setLogin] = useState(false);
  const [room, setRoom] = useState(null);
  const [name, setName] = useState("");
  const [cloudInfo, setCloudInfo] = useState(null);
  const [errorMsg, setErrorMsg] = useState("");

  const CreateNewRoom = async () => {
    if (name === "") {
      setErrorMsg("Please enter a nickname")
      return
    }

    console.log("Creating room...")
    const response = await fetch(CREATENEWROOMAPI, {
      method: 'POST'
    })
    const data = await response.json()
    setCloudInfo(data)
    setRoom(data.room_id)
    setLogin(true)
  }

  const JoinRoom = async () => {
    if (room === null) {
      setErrorMsg("Please enter a room pin")
      return
    }
    if (name === "") {
      setErrorMsg("Please enter a nickname")
      return
    }

    console.log("joining room")

    const response = await fetch(CHECKROOMAPI + room)
    const data = await response.json()
    
    if (data.result) {
      console.log("room exist")

      const response = await fetch(GETROOMAPI + room)
      const data = await response.json()
      console.log(data)
      setCloudInfo(data.result)
      setLogin(true)
    } else {
      console.log("room does not exist")
    }
    
  }

  useEffect(() => {

  }, [isLogin])

  return (
    <div style={{height: '100vh'}}>
      { 
        isLogin ? <CodeWindow room_id={room} cloudInfo={cloudInfo} name={name}/> :
        (
          <div className='Container' style={{background: '#C5D4D2', height: '100vh'}}>
            <div className='row'>
              <h1 className='my-5' style={{color: '#473E3E'}}>CodeBooks</h1>
            </div>
            <div className='row'>
              <div className='col' />
              <div className='col' >
                <div className='p-4 mt-5' style={{background: '#17C3A5', borderRadius: 10}}>
                  <Form.Control
                    size="lg"
                    id="pin"
                    className='my-2'
                    placeholder='Room PIN'
                    onChange={(e) => setRoom(e.target.value)}
                  />
                  <Form.Control
                    id="name"
                    size="lg"
                    className='my-2'
                    placeholder='Nickname'
                    onChange={(e) => setName(e.target.value)}
                  />
                  <Form.Text className="text-muted">
                    {errorMsg}
                  </Form.Text>
                  <div className='row mb-2 mt-5'>
                    <div className='col'>
                      <Button size="lg" style={{background: '#473E3E', color: '#FCFCFC', width: '100%'}} onClick={() => JoinRoom()}>Enter</Button>
                    </div>
                  </div>
                  <div className='row'>
                      <div className='col'>
                        <Button size="lg" style={{background: '#473E3E', color: '#FCFCFC', width: '100%'}} onClick={() => CreateNewRoom()}>Create New Room</Button>
                      </div>
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