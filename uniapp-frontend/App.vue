<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'

onLaunch(() => {
  console.log('App Launch')
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
@import 'uview-ui/index.scss';

page {
  background-color: #F5F5F5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

/* 主题色 */
$primary-color: #00DC82;
$text-primary: #333333;
$text-secondary: #666666;
$text-placeholder: #999999;
$border-color: #EEEEEE;
$background-color: #F5F5F5;
$card-bg: #FFFFFF;

/* 涨跌颜色 */
$up-color: #00DC82;
$down-color: #FF4D4D;
</style>
