<template>
  <v-app>
    <Header :keycloak="keycloak"/>

    <v-main>
      <v-container
          class="d-flex justify-center mb-6"
          fluid>
        <router-view/>
      </v-container>
    </v-main>

  </v-app>
</template>

<script>
import Header from './components/app/Header'

export default {
  name: 'App',

  components: {
    Header,
  },

  props: {
    keycloak: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    //
  }),
};
</script>

<style>
.theme--light.v-application {
  background-image: linear-gradient(#e6f0f7, #e6f0f0) !important;
}

html {
  font-size: 1rem;
}

 .v-data-table > .v-data-table__wrapper > table > thead > tr > th {
   font-size: 1rem !important;
   font-weight: lighter;
 }

.v-data-table > .v-data-table__wrapper > table > tbody > tr > td {
  font-size: 1rem !important;
}
</style>
