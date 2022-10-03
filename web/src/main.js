import Vue from 'vue'
import App from './App.vue'

// import ElementUI from 'element-ui'
// import 'element-ui/lib/theme-chalk/index.css'
// Vue.use(ElementUI)

import { Button, Input, Image, Icon, Row, Col, Card, Pagination, Link, Backtop } from 'element-ui'

Vue.use(Button)
Vue.use(Input)
Vue.use(Image)
Vue.use(Icon)
Vue.use(Row)
Vue.use(Col)
Vue.use(Card)
Vue.use(Pagination)
Vue.use(Link)
Vue.use(Backtop)

Vue.config.productionTip = false

new Vue({
  el: '#app',
  render: h => h(App)
})
