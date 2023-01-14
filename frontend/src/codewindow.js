import React, { useState } from 'react'


import Button from 'react-bootstrap/Button';

import { highlight, languages } from 'prismjs/components/prism-core';
import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-javascript';
import 'prismjs/themes/prism.css'; //Example style, you can use another
import Editor from 'react-simple-code-editor'

const sampleCode = `
#include <stdio.h>

int main() {
  printf("Hello World!");
  return 0;
}
`

const SUBMITENDPOINT = "http://localhost:8080/api/v1/judge/submit"

function CodeWindow() {
  const [code, setCode] = useState(sampleCode)

  const submitCode = async () => {
    const source_code = btoa(code)
    const bodyData = {
      source_code: source_code,
      language: "c",
      room_id: "1234",
    }
    const response = await fetch(SUBMITENDPOINT, {
      method: 'POST',
      body: JSON.stringify(bodyData)
    })
    const data = await response.json()
    console.log(data)
  }

  return (
    <div className='Container'>
      <div className='row'>
        <div className='col'>
          <h1>This is the code window</h1>
          <div style={{backgroundColor: '#e6ddc8'}}>
            <Editor 
              value={code}
              onValueChange={code => setCode(code)}
              highlight={code => highlight(code, languages.js)}
              padding={10}
              style={{
                fontFamily: '"Fira code", "Fira Mono", monospace',
                fontSize: 14
              }}
              
            />
          </div>
          <Button onClick={() => submitCode()}>RUN</Button>
        </div>
      </div>
    </div>
  )
}

export default CodeWindow