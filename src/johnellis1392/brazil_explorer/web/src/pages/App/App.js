import React, { Component } from 'react'
import PropTypes from 'prop-types'
import styles from './App.scss'

export class App extends Component {

  constructor() {
    super()
    this.state = {}
  }

  static propTypes = {
    /* */
  }

  static defaultProps = {
    /* */
  }

  render() {
    return (
      <div className={styles.wrapper}>
        Hello, World!
      </div>
    )
  }

}

export default App
