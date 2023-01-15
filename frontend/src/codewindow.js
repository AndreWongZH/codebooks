
import React, { useEffect, useState } from 'react'
import { encode, decode } from 'js-base64';


import Button from 'react-bootstrap/Button';

import { highlight, languages } from 'prismjs/components/prism-core';
import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-javascript';
import 'prismjs/components/prism-c';
import 'prismjs/components/prism-cpp';
import 'prismjs/components/prism-go';
import 'prismjs/components/prism-python';
import 'prismjs/themes/prism.css';
import Editor from 'react-simple-code-editor'
import io from 'socket.io-client';

import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Dropdown from 'react-bootstrap/Dropdown';
import DropdownButton from 'react-bootstrap/DropdownButton';

import 'prismjs/components/prism-python'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'

const langMap = {
  "c": languages.c,
  "c++": languages.cpp,
  "golang": languages.go,
  "python3": languages.python
}

const langMap2 = {
  "c": "c",
  "c++": "cpp",
  "golang": "go",
  "python3": "python"
}

const sampleCode = `
// this is a c code

#include <stdio.h>

int main() {
  printf("Hello World!");
  return 0;
}
`
const hostname = 'http://52.221.121.128:8080'
// const hostname = 'http://localhost:8080'

const SUBMITENDPOINT = hostname + "/api/v1/judge/submit"
const socket = io(hostname);

function CodeWindow({room_id, cloudInfo, name}) {
  const [code, setCode] = useState(sampleCode)
  const [lang, setLang] = useState("c")
  const [output, setOutput] = useState("~/Desktop/" + name + "$")
  const [users, setUsers] = useState([name])
  const [lock, setLock] = useState(true)

  useEffect(() => {
    if (cloudInfo) {
      console.log(cloudInfo)
      setCode(cloudInfo.source_code)
      setLang(cloudInfo.language)
    }
  }, [cloudInfo])

  const updateCode = (code) => {
    if (!lock) {
      return
    }
    setCode(code)
    const obj = {
      "source_code": code,
      "user": name,
      "room_id": room_id,
    }
    socket.emit("edit", JSON.stringify(obj))
    
  }

  const changeColor = async (userid) => {
    const el = document.getElementById(userid);
    if (el) {
      el.style.backgroundColor = 'red'
      await new Promise(r => setTimeout(r, 500));
      el.style.backgroundColor = '#31a346'
    }
  }

  const changeLock = async () => {
    setLock(false)
    await new Promise(r => setTimeout(r, 500));
    setLock(true)
  }

  useEffect(() => {

    socket.on("newcode", (msg) => {
      const obj = JSON.parse(msg)
      if (obj.user !== name) {
        changeColor(obj.user)
        changeLock()
        setCode(obj.source_code)
      }
      console.log(msg)
    })

    socket.emit("joinroom", JSON.stringify({
      "room_id": room_id,
      "source_code": "",
      "user": name,
    }))

    socket.on("ping", () => 
      socket.emit("pongpong", JSON.stringify({
        "room_id": room_id,
        "source_code": "",
        "user": name,
      }))
    )

    socket.on("active_users", (msg) => {
      console.log(msg)
      const set = new Set(msg);
      set.delete(name)
      setUsers([name, ...Array.from(set)])
    })

    return () => {
      socket.off('newcode');
    };
  }, []);

  const submitCode = async () => {
    setOutput("Executing code ... .. .")
    const source_code = encode(code)
    const bodyData = {
      source_code: source_code,
      language: lang,
      room_id: "1234",
    }
    console.log(bodyData)
    const response = await fetch(SUBMITENDPOINT, {
      method: 'POST',
      body: JSON.stringify(bodyData)
    })
    const data = await response.json()
    console.log(data)
    let outputData
    if (data.compile_output === "") {
      if (data.stderr) {
        outputData = decode(data.stderr)
      } else {
        if (data.stdout) {
          outputData = "> \n" + decode(data.stdout) + " \nTime taken: " + data.time + "s"
        } else {
          outputData = "> \n"
        }
      }
    } else {
      outputData = decode(data.compile_output)
    }
    
    setOutput(outputData)
  }

  return (
    <div className='Container'>

      <div className='row'>
        <div className='col-9'>
          <div style={{backgroundColor: '#e6ddc8', height: '70vh', overflow: 'scroll'}}>
            <Editor
              value={code}
              onValueChange={updateCode}
              highlight={code => highlight(code, langMap[lang], langMap2[lang])}
              padding={10}
              style={{
                fontFamily: '"Fira code", "Fira Mono", monospace',
                fontSize: 14,
              }}
            />
          </div>
        </div>
        <div className='col-3'>
          <h3 className='mt-2'>Room</h3>
          <div className='py-2' style={{background: '#D9D9D9', borderRadius: '20px', fontSize: '50px'}}>{room_id}</div>

          <h3 className='mt-5'>Users</h3>
          <div style={{width: '100%', margin: 'auto'}}>
            {
              users.map((name, index) => {
                return (
                  <div key={name + "container"} className="py-2 mb-2" style={{display: 'flex', justifyContent: 'start', background: '#D9D9D9', alignItems: 'center', borderRadius: '20px'}}>
                      <div key={name + "color"} id={name} className="ms-4" style={{borderRadius: '100%', background: '#31a346', width: '30px', height: '30px'}} />
                      <div key={name + "name"} className="ms-4" >{name}</div>
                  </div>
                )
              })
            }
          </div>
        </div>
      </div>
      

      <div className='row' style={{background: "#CCCCCC"}}>
        <div className='col'>
          <Button size="lg" style={{width: '100%', borderRadius: '0px', background: '#70B6DD', borderColor: '#70B6DD'}} onClick={() => submitCode()}>RUN</Button>
        </div>
        <div className='col-8'></div>
        <div className='col'>
          {
            <DropdownButton
              style={{width: '100%', borderRadius: '0px', background: '#70B6DD', borderColor: '#70B6DD', color: '#EEEEEE'}}
              as={ButtonGroup}
              size="lg"
              variant='#70B6DD'
              // key={variant}
              // id={`dropdown-variants-${variant}`}
              // variant={variant.toLowerCase()}
              title={lang}
              
            >
              <Dropdown.Item eventKey="c" onClick={(e) => setLang(e.target.text)} >c</Dropdown.Item>
              <Dropdown.Item eventKey="c++" onClick={(e) => setLang(e.target.text)} >c++</Dropdown.Item>
              <Dropdown.Item eventKey="golang" onClick={(e) => setLang(e.target.text)} >golang</Dropdown.Item>
              <Dropdown.Item eventKey="python3" onClick={(e) => setLang(e.target.text)} >python3</Dropdown.Item>
            </DropdownButton>
          }
          {/* <Form.Select
            aria-label="Default select example"
            onChange={(e) => setLang(e.target.value)}
          >
            <option value="c">c</option>
            <option value="c++">c++</option>
            <option value="golang">golang</option>
            <option value="python3">python3</option>
          </Form.Select> */}
        </div>
      </div>

      <div className='row'>
        <div className='col'>
          <div className="ps-2" style={{color: '#a3c2ab', background: '#1a211c', height: '23vh', whiteSpace: "break-spaces", textAlign: "left", overflow: 'scroll'}}>
            {output}
          </div>
        </div>
      </div>
    </div>
  )
}

export default CodeWindow