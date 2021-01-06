<template>
  <v-select
      :value="value"
      :items="users"
      :label="label ? label : 'Users'"
      :item-text="getUserFullName"
      item-value="id"
      multiple
      dense
      @change="$emit('change', '')"
      @input="emitChange"
  >
    <template v-slot:prepend-item>
      <v-list-item
          ripple
          @click="toggle"
          v-if="users && users.length > 0"
      >
        <v-list-item-action>
          <v-icon :color="value !== undefined && value.length > 0 ? 'indigo darken-4' : ''">
            {{ icon }}
          </v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>
            Select All
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider
          class="mt-2"
          v-if="users && users.length > 0"
      ></v-divider>
    </template>

    <template v-slot:selection="{item, index}">
      <slot name="selection" :item="item" :fullName="getUserFullName(item)" :index="index">
        {{ getUserFullName(item) }}
        <span v-if="value && index < (value.length - 1)" class="pr-1">,</span>
      </slot>
    </template>
  </v-select>
</template>

<script>
export default {
  name: "UserSelect",
  props: {
    value: {
      type: Array,
      required: false,
    },
    users: {
      type: Array,
      required: true,
    },
    label: {
      type: String,
      required: false,
    },
  },
  methods: {
    getUserFullName(user) {
      return `${user.firstName} ${user.lastName}`
    },
    toggle() {
      this.$nextTick(() => {
        if (this.value && this.value.length === this.users.length) {
          this.$emit('input', [])
          this.$emit('change', '')
        } else {
          this.$emit('input', this.users.slice().map(u => u.id))
          this.$emit('change', '')
        }
      })
    },
    emitChange(selectedUsers) {
      this.$emit('input', selectedUsers)
    },
  },
  computed: {
    icon() {
      if (this.value) {
        if (this.value.length === this.users.length) return 'mdi-close-box'
        if (this.value.length > 0) return 'mdi-plus-box'
      }
      return 'mdi-checkbox-blank-outline'
    },
  },
}
</script>