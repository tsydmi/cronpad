<template>
  <tbody>
  <tr
      v-for="item in integrations"
      :key="item.name"
  >
    <td>
      <span v-bind:class=" item.enabled ? 'success--text' : 'secondary--text'">
        {{ item.name }}
      </span>
    </td>
    <td class="d-flex justify-end">
      <v-btn
          class="ma-2"
          small
          outlined
          color="indigo"
          :disabled="!item.enabled"
          v-if="item.hasConfig"
      >
        <v-icon>mdi-pencil</v-icon>
      </v-btn>
      <v-btn
          class="ma-2"
          small
          color="success"
          outlined
          :disabled="!item.enabled"
          v-if="!item.hasConfig"
      >
        add
      </v-btn>
      <v-btn
          class="ma-2"
          small
          color="error"
          outlined
          :disabled="!item.enabled"
          v-if="item.hasConfig"
      >
        remove
      </v-btn>
    </td>
  </tr>
  </tbody>
</template>

<script>
export default {
  name: "IntegrationTable",

  data: () => ({
    integrations: [
      {
        name: "Google Calendar",
        enabled: true,
        disabled: false
      },
      {
        name: "Bitbucket",
        enabled: false,
        disabled: true
      },
      {
        name: "GitLab",
        enabled: false,
        disabled: true
      }
    ],
  }),
}
</script>

<style scoped>

</style>