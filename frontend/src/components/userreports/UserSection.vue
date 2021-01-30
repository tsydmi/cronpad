<template>
  <div>
    <v-row>
      <v-col cols="8">
        <v-select
            v-model="selectedProject"
            :items="projects"
            @change="refreshUsers"
            label="PROJECT"
            item-value="id"
            item-text="name"
            dense
            clearable
        ></v-select>
      </v-col>
      <v-col>
        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-btn
                icon
                @click="clearFilters"
                v-bind="attrs"
                v-on="on"
                :disabled="!selectedProject"
            >
              <v-icon>far fa-times-circle</v-icon>
            </v-btn>
          </template>
          <span>Clear filters</span>
        </v-tooltip>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-list flat v-if="users && users.length > 0">
          <v-list-item-group
              color="primary"
              @change="emitChange"
              v-model="index"
          >
            <v-list-item
                v-for="(user, i) in users"
                :key="i"
            >
              <v-list-item-content
              >
                <v-list-item-title> {{ user.firstName }} {{ user.lastName }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>
        <span class="ma-5" v-else>
          No users found
        </span>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import UserService from "@/service/UserService"
import ProjectService from "@/service/ProjectService"

export default {
  name: "UserSection",

  props: {
    value: {
      type: Object,
      required: false,
    },
  },
  data: () => ({
    index: null,
    selectedProject: null,

    users: [],
    projects: [],
  }),
  methods: {
    emitChange(item) {
      this.$emit('change', this.getUserByIndex(item))
    },
    clearFilters() {
      this.selectedProject = null
      this.refreshUsers()
    },
    refreshUsers() {
      if (this.selectedProject) {
        ProjectService.getUsers(this.selectedProject)
            .then((response) => {
              this.users = response.data
              this.updateIndex()
            })
      } else {
        UserService.findAll()
            .then((response) => {
              this.users = response.data
              this.updateIndex()
            })
      }
    },
    refreshProjects() {
      ProjectService.findAll(true)
          .then((response) => this.projects = response.data)
    },
    getUserByIndex(i) {
      if (i === undefined || i === null) {
        return null
      } else {
        return this.users[i]
      }
    },
    updateIndex() {
      if (this.value) {
        this.index = this.users
            .findIndex(u => u.id === this.value.id);
      } else {
        this.index = undefined
      }
    },
  },
  created() {
    this.refreshUsers()
    this.refreshProjects()
  }
}
</script>
