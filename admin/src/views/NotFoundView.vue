<template>
  <div class="adm-page">
    <el-result icon="warning" title="页面不存在" sub-title="这个后台页面不存在或已经被删除。">
      <template #extra>
        <el-button type="primary" @click="$router.push('/pages')">返回页面列表</el-button>
      </template>
    </el-result>
  </div>
</template>
