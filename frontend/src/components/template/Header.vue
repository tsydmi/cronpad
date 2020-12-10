<template>
  <v-app-bar
      color="white"
      app
  >
    <div class="d-flex align-center">
      <v-img
          alt="Cronpad Logo"
          class="shrink mr-2"
          contain
          src="@/assets/icon.svg"
          transition="scale-transition"
          width="40"
      />

      <h1 class="display-1 font-weight-bold mb-3 primary--text pt-2">
        Cronpad
      </h1>
    </div>

    <v-toolbar>
      <v-toolbar-items>
        <v-divider vertical></v-divider>

        <v-btn text color="primary" to="/">
          <i class="fas fa-home"></i>
        </v-btn>

        <v-divider vertical></v-divider>

        <v-btn text color="primary" to="/working-time">
          <i class="far fa-calendar-alt pr-2"></i>
          Working time
        </v-btn>

        <v-divider vertical></v-divider>

        <v-btn text color="secondary" to="/settings">
          <i class="fas fa-cogs"></i>
        </v-btn>

        <v-divider vertical></v-divider>
      </v-toolbar-items>
    </v-toolbar>


    <div class="d-flex align-center">
      <h4 class="pa-2 secondary--text">{{ user.fullName }}</h4>
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
                color="purple"
            >
              <span class="white--text">{{ user.initials }}</span>
            </v-avatar>
          </v-btn>
        </template>
        <v-card>
          <v-list-item-content class="justify-center">
            <div class="mx-auto text-center">
              <v-avatar
                  color="purple"
              >
                <span class="white--text headline">{{ user.initials }}</span>
              </v-avatar>
              <p class="caption mt-1">
                {{ user.email }}
              </p>
              <v-divider class="my-3"></v-divider>
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
export default {
  props: {
    keycloak: {
      type: Object,
      required: true,
    }
  },
  data: () => ({
    user: {
      initials: '',
      fullName: '',
      email: '',
    },
  }),
  methods: {
    getInitials(fullName) {
      var result = ""

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
    this.keycloak.loadUserInfo()
        .then(() => {
          this.user.initials = this.getInitials(this.keycloak.userInfo.name)
          this.user.fullName = this.keycloak.userInfo.name
          this.user.email = this.keycloak.userInfo.email
        })
  }
}
</script>

<style>
.v-sheet.v-toolbar:not(.v-sheet--outlined) {
  box-shadow: none !important
}
</style>