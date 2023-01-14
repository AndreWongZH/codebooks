
import React, { useEffect, useState } from 'react'
import { encode, decode } from 'js-base64';


import Button from 'react-bootstrap/Button';


import { highlight } from 'prismjs/components/prism-core';
import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-javascript';
import 'prismjs/themes/prism.css';
import Editor from 'react-simple-code-editor'
import Form from 'react-bootstrap/Form';
import io from 'socket.io-client';

import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Dropdown from 'react-bootstrap/Dropdown';
import DropdownButton from 'react-bootstrap/DropdownButton';

import Prism from 'prismjs';
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'

const sampleCode = `
// this is a c code

#include <stdio.h>

int main() {
  printf("Hello World!");
  return 0;
}
`

const SUBMITENDPOINT = "http://localhost:8080/api/v1/judge/submit"
const socket = io("http://127.0.0.1:8080/");

function CodeWindow({room_id, cloudInfo, name}) {
  const [code, setCode] = useState(sampleCode)
  const [lang, setLang] = useState("c")
  const [output, setOutput] = useState("~/Desktop/" + name + "$")

  useEffect(() => {
    if (cloudInfo) {
      console.log(cloudInfo)
      setCode(cloudInfo.source_code)
      setLang(cloudInfo.language)
    }
  }, [cloudInfo])

  const updateCode = (code) => {
    setCode(code)
    const obj = {
      "source_code": code,
      "user": name,
      "room_id": room_id,
    }
    socket.emit("edit", JSON.stringify(obj))
    
  }

  useEffect(() => {
    socket.on("noedit", () => {
      console.log("noedit")
    })

    socket.on("edit", () => {
      console.log("edit")
    })

    socket.on("newcode", (msg) => {
      const obj = JSON.parse(msg)
      if (obj.user !== name) {
        setCode(obj.source_code)
      }
      console.log(msg)
    })

    return () => {
      socket.off('noedit');
      socket.off('edit');
    };
  }, []);

  const submitCode = async () => {
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

    let outputData
    if (data.compile_output === "") {
      outputData = "> \n" + decode(data.stdout) + " \nTime taken: " + data.time + "s"
    } else {
      console.log(data.compile_output)
      console.log()
      let errorMsg = data.compile_output.split('\n');
      const errorMsgDecoded = errorMsg.map(el => {
        return decode(el)
      });
      const errorMsgString = errorMsgDecoded.join('\n')
      console.log(errorMsgString)
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
              highlight={code => highlight(code, Prism.languages.cpp, 'cpp')}
              padding={10}
              style={{
                fontFamily: '"Fira code", "Fira Mono", monospace',
                fontSize: 14,
              }}
            />
          </div>
        </div>
        <div className='col-3'>
          <h3>Room</h3>
          <h6>{room_id}</h6>
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