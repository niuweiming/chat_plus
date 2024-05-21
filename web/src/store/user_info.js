import {defineStore} from 'pinia';
import { cellGroupProps } from 'vant';

export const user_infoStore = defineStore('user_infos',{
    state: () => {
      return {
        username: '',
       };
    },
    actions: {
        getusername(){
        return this.username;
        },
        setusername(username_pay){
            console.log('设置成功')
          this.username = username_pay;
        }
    },})
      
   