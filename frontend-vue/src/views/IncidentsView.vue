<script setup lang="ts">
import { createIncident, isAuthenticated, loadIncidents, logout } from '@/api';
import AddIncidentForm from '@/components/AddIncidentForm.vue';
import AppHeader from '@/components/AppHeader.vue';
import IncidentList from '@/components/IncidentList.vue';
import type { Severity, CreateIncidentRequest, Incident } from '@/types';
import { onMounted, ref } from 'vue';
const incidents = ref<Incident[]>([])

async function handleLogout() {
  await logout()
  window.location.href = "/"
}

onMounted(async() => {
  if (await isAuthenticated() == false) {
    window.location.href = "/"
  }
  incidents.value = await loadIncidents()
})

</script>

<template>
  <main>
    <AppHeader @submit="handleLogout"></AppHeader>
    <div class="page">
      <div class="dash-head">
        <div>
          <p class="eyebrow">Active board</p>
          <h1 class="page-title">Incidents</h1>
        </div>
        <div class="spacer"></div>
        <RouterLink :to="{name: 'incidents-new'}" class="btn btn-primary">+ New Incident</RouterLink>
      </div>
    </div>
  </main>
</template>

<style scoped>
.dash-head {
  align-items: flex-end;
  display: flex;
}
</style>