import React, { Component } from 'react'
import ReactDOM from 'react-dom'
import 'jquery'

import { Editor } from 'pages'


const data = {
  root: {
    name: "Node 1",
  },
}


const App = ({}) => {
  return (
    <div>
      Hello, World! From React
      <Editor data={data}/>
    </div>
  )
}

ReactDOM.render(
  <App/>,
  document.getElementById('root')
)
