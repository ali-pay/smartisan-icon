<template>
  <div id="app" class="wrap">
    <el-input clearable placeholder="请输入内容" v-model="inputName">
      <el-button slot="append" icon="el-icon-search">查找</el-button>
    </el-input>
    <el-row :gutter="20" v-if="showUrls">
      <el-col :span="6" v-for="(o, index) in pageSize" :key="index">
        <el-card shadow="hover" :body-style="{ padding: '0px' }" v-if="findUrl(index)">
          <el-image lazy :src="findUrl(index)"></el-image>
          <span class="name">{{ findName(findUrl(index)) }}</span>
          <span class="download" @click="downloadImage(findUrl(index), findName(findUrl(index)))"><i class="el-icon-download"></i></span>
        </el-card>
      </el-col>
    </el-row>
    <el-pagination
      v-if="showUrls"
      background
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page.sync="currentPage"
      :page-sizes="[12, 24, 48, 96, 999]"
      :page-size="pageSize"
      :total="showUrls.length"
      layout="total, sizes, prev, pager, next, jumper"
    >
    </el-pagination>
  </div>
</template>

<script>
import iconUrls from '../../downloader/icon_url.json'
export default {
  name: 'App',
  data() {
    return {
      iconUrls,
      showUrls: null,
      currentPage: 1,
      pageSize: 12,
      inputName: '',
      timer: null
    }
  },
  watch: {
    inputName(name) {
      if (name === '') {
        this.showUrls = this.iconUrls
        return
      }
      this.debounce(() => {
        this.showUrls = this.iconUrls.filter(item => item.indexOf(name, 35) !== -1)
      })
    }
  },
  mounted() {
    // 按名称排序
    this.iconUrls = this.iconUrls.sort((a, b) => a.localeCompare(b))
    this.showUrls = this.iconUrls
  },
  methods: {
    handleSizeChange(val) {
      this.pageSize = val
    },
    handleCurrentChange() {
      // 滚动到页面顶部
      let top = document.documentElement.scrollTop || document.body.scrollTop
      let timer = setInterval(() => {
        document.body.scrollTop = document.documentElement.scrollTop = top -= 50
        if (top <= 0) clearInterval(timer)
      }, 10)
    },
    findUrl(index) {
      return this.showUrls[index + (this.currentPage - 1) * this.pageSize]
    },
    findName(url) {
      return url.substring(35, url.lastIndexOf('/'))
    },
    downloadImage(src, name) {
      let image = new Image()
      image.setAttribute('crossOrigin', 'anonymous')
      image.onload = () => {
        let canvas = document.createElement('canvas')
        canvas.width = image.width
        canvas.height = image.height
        let context = canvas.getContext('2d')
        context.drawImage(image, 0, 0, image.width, image.height)
        let url = canvas.toDataURL('image/png')
        let a = document.createElement('a')
        let event = new MouseEvent('click')
        a.download = name + '.png'
        a.href = url
        a.dispatchEvent(event)
      }
      image.src = src
    },
    debounce(fn, wait = 500) {
      if (this.timer) clearTimeout(this.timer)
      this.timer = setTimeout(() => {
        fn()
      }, wait)
    }
  }
}
</script>

<style>
html {
  background: #fafafb;
}
.wrap {
  width: 1220px;
  margin: 0 auto;
}
.name,
.download {
  width: 100%;
  height: 40px;
  line-height: 40px;
  display: block;
  color: #666;
  font-size: 12px;
  text-align: center;
}
.download {
  display: none;
  font-size: 22px;
  color: #fff;
  background: rgba(13, 10, 49, 0.9);
  cursor: pointer;
}
.download:hover {
  color: red;
  background: #0d0a31;
}
.el-card {
  width: 252px;
  height: 278px;
  border-radius: 8px;
  margin: 20px 0;
}
.el-card:hover .name {
  display: none;
}
.el-card:hover .download {
  display: block;
}
.el-image {
  display: block;
  margin: 0 auto;
  margin-top: 18px;
  width: 220px;
  height: 220px;
  background: #f7f7f7;
}
.is-hover-shadow:hover {
  box-shadow: 1px 1px 10px 2px #ccc !important;
}
.el-input {
  width: 401px;
}
</style>
