<template>
  <v-app-bar
      color="white"
      app
  >
    <div class="d-flex align-center">
      <v-img
          alt="Cronpad Logo"
          class="shrink ma-5"
          contain
          src="@/assets/icon.svg"
          transition="scale-transition"
          width="40"
      />

      <span> CRONPAD </span>
    </div>

    <v-container>
      <v-tabs>
        <v-tabs-slider color="yellow"></v-tabs-slider>

        <v-tab to="/">
          <em class="fas fa-home"></em>
        </v-tab>
        <v-tab to="/working-time">
          Working time
        </v-tab>
        <v-tab to="/projects" v-if="hasAdminRole">
          Projects
        </v-tab>
        <v-tab to="/tags" v-if="hasManagerRole">
          Tags
        </v-tab>
        <v-tab to="/user-reports" v-if="hasAdminRole">
          User Reports
        </v-tab>
        <v-tab to="/project-statistics" v-if="hasManagerRole">
          Project Statistics
        </v-tab>
        <v-tab to="/settings">
          <em class="fas fa-cogs"></em>
        </v-tab>

      </v-tabs>
    </v-container>

    <div class="d-flex align-center">
      <h4 class="pa-2 inactive--text">{{ user.fullName }}</h4>
      <v-menu
          rounded
          offset-y
      >
        <template v-slot:activator="{ on }">
          <v-btn
              icon
              v-on="on"
          >
            <v-avatar
                color="primary"
            >
              <span class="white--text">{{ user.initials }}</span>
            </v-avatar>
          </v-btn>
        </template>
        <v-card>
          <v-list-item-content class="justify-center">
            <div class="mx-auto text-center">
              <v-btn
                  disabled
                  depressed
                  rounded
                  text
              >
                Edit profile
              </v-btn>
              <v-divider class="my-3"></v-divider>
              <v-btn
                  depressed
                  rounded
                  text
                  @click="signOut"
              >
                Sign out
              </v-btn>
            </div>
          </v-list-item-content>
        </v-card>
      </v-menu>
    </div>
  </v-app-bar>
</template>

<script>
const ADMIN_ROLE = 'admin'
const MANAGER_ROLE = 'project-manager'

export default {
  props: {},
  data: () => ({
    user: {
      initials: '',
      fullName: '',
      email: '',
    },
    hasAdminRole: false,
    hasManagerRole: false,
  }),
  methods: {
    getInitials(fullName) {
      let result = ""

      fullName.split(" ").forEach(name => {
        if (name.trim()) {
          result += name[0].toUpperCase()
        }
      })

      return result
    },
    signOut() {
      this.keycloak.logout()
    }
  },
  created() {
    this.$keycloak.loadUserInfo()
        .then(() => {
          this.user.initials = this.getInitials(this.$keycloak.userInfo.name)
          this.user.fullName = this.$keycloak.userInfo.name
          this.user.email = this.$keycloak.userInfo.email

          this.hasAdminRole = this.$keycloak.hasRealmRole(ADMIN_ROLE)
          this.hasManagerRole = this.$keycloak.hasRealmRole(MANAGER_ROLE)
        })
  }
}
</script>

<style>
.v-sheet.v-toolbar:not(.v-sheet--outlined) {
  box-shadow: none !important
}
</style>