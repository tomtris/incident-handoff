<template>
  <span class="badge" :class="'sev-' + level">{{ level }}</span>
</template>

<script>
export default {
  name: "SeverityBadge",
  props: { level: { type: String, default: "SEV3" } },
};
</script>