---
# 仅用于语言自动跳转
title: Redirect
layout: page
head:
  - - meta
    - name: robots
      content: noindex, nofollow
---

<script setup>
if (typeof window !== 'undefined') {
  const lang = (navigator.language || navigator.userLanguage || '').toLowerCase()
  const target = lang.startsWith('zh') ? '/autocft/zh/' : '/autocft/en/'
  window.location.replace(target)
}
</script>

<div style="padding:2rem 0;font-size:14px;">
  <p>Redirecting by locale...</p>
  <p>
    If not redirected automatically:
    <a href="/en/">English</a> | <a href="/zh/">中文</a>
  </p>
  <noscript>
    JavaScript disabled. Please choose: <a href="/en/">English</a> | <a href="/zh/">中文</a>
  </noscript>
</div>
