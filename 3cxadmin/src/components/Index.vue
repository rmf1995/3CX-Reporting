<template>
   <div>
      <v-card>
         <v-card-title>
            3CX Servers
            <v-spacer></v-spacer>
            <v-text-field
               v-model="search"
               append-icon="mdi-magnify"
               label="Search"
               single-line
               hide-details
               ></v-text-field>
         </v-card-title>
         <v-data-table
            :headers="headers"
            :items="items"
            :search="search"
            sort-by="calories"
            class="elevation-1"
            >
            <template v-slot:item.actions="{ item }">
               <v-btn
                  class="ma-2"
                  color="warning"
                  dark
                  :to="{name: 'Edit', params: { id: item.id }}"
                  >
                  <v-icon>
                     mdi-pencil
                  </v-icon>
               </v-btn>
               <v-btn
                  class="ma-2"
                  color="error"
                  dark
                  @click="deleteItem(item.id)"
                  >
                  <v-icon>
                     mdi-delete
                  </v-icon>
               </v-btn>
            </template>
         </v-data-table>
      </v-card>
   </div>
</template>


<script>
    export default {
        data(){
            return{
                search: '',
                headers: [
                {
                    text: 'VM Name',
                    align: 'start',
                    filterable: true,
                    value: 'Name',
                },
                { text: 'VM Location', value: 'Location' },
                { text: 'PwState ID', value: 'PwStateID' },
                { text: 'Actions', value: 'actions', filterable: false, sortable: false },
                ],
                items: []
            }
        },

        created: function()
        {
            this.fetchItems();
        },

        methods: {
            fetchItems()
            {
              let uri = 'https://URL/v1/getAll3CX/';
              this.axios.get(uri).then((response) => {
                  this.items = response.data;
              });
            },
            deleteItem(id)
            {
                if(confirm("Do you really want to delete?")){
                    let uri = 'https://URL/v1/3cx/delete/'+id;
                    this.items.splice(id, 1);
                    this.axios.delete(uri).then((response) => { this.$router.push({name: 'Index'}); });
                }
            },
        }
    }
</script>