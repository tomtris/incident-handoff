<script setup lang="ts">
// async function handleCreateIncident({title, service, severity}: {title: string, service: string, severity: string}) {
//   function isSeverity(v: string) : v is Severity {
//     return v === "SEV1" || v === "SEV2" || v === "SEV3"
//   }

//   if (!isSeverity(severity)) {
//     createIncidentFormError.value = "Severity must be SEV1, SEV2 or SEV3"
//     return
//   }

//   const incReq : CreateIncidentRequest = {
//     title: title,
//     service: service,
//     severity: severity
//   }

//   try {
//     const res = await createIncident(incReq)
//     createIncidentFormError.value = ''
//     incidents.value.push(res)
//   } catch (e) {
//     createIncidentFormError.value = (e as Error).message
//   }
// }
</script>

<template>

</template>