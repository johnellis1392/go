import React, { Component } from 'react'
import PropTypes from 'prop-types'
import styles from './editor.scss'


export class Node extends Component {

  constructor() {
    super()
    this.state = {}
  }

  static propTypes = {
    name: PropTypes.string.isRequired,
  }

  static defaultProps = {}

  render() {
    return (
      <div>
        <span>* {this.props.name}</span>
      </div>
    )
  }


}

export class Editor extends Component {

  constructor() {
    super()
    this.state = {}
  }

  static propTypes = {
    data: PropTypes.shape({
      root: PropTypes.object.isRequired,
    }).isRequired,
  }

  static defaultProps = {}

  render() {
    const nodeProps = Object.assign({}, {
      ...this.props.data.root,
    })

    return (
      <div>
        Editor
        <Node {...nodeProps}/>
      </div>
    )
  }

}

export default Editor
