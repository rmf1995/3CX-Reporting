<template>
  <v-app id="inspire">
    <v-app-bar
      app
      flat
    >
      <v-container class="py-0 fill-height">
        <v-btn
          v-for="item in menuItems"
          :key="item.title"
          :to="item.path"
          text
        >
          {{ item.title }}
        </v-btn>



      </v-container>
    </v-app-bar>

    <v-main >
      <v-container>
        <v-row>


          <v-col>
            <v-sheet
              min-height="70vh"
              rounded="lg"
            >
              <router-view></router-view>
            </v-sheet>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  export default {
    data: () => ({
      menuItems: [
          { title: 'All 3CX', path: '/index' },
          { title: 'Add 3CX', path: '/create' }
     ],
    }),
  }
</script>