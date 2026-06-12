<script setup lang="ts">
import { ref } from 'vue';
import LoginForm from '../components/LoginForm.vue'
import { login } from '@/api.ts';

const error = ref('')
async function handleLogin(payload : {username:string; password:string})  {
    error.value = ''
    try {
      await login(payload.username, payload.password)
      window.location.href = "/incidents"
    } catch (e) {
      error.value = (e as Error).message ?? 'Authentication failed'
    }
}
</script>

<template>
  <main>
    <LoginForm :error="error" @submit="handleLogin"/>
  </main>
</template>
