import React from 'react'
import ReactDOM from 'react-dom'
import { IntlProvider } from 'react-intl'
import { App } from './pages'
import 'jquery'

const locale = 'en'

ReactDOM.render(
  <IntlProvider locale={locale}>
    <App/>
  </IntlProvider>,
  document.getElementById("root")
)
