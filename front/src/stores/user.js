import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', () => {
    const userId = ref('');

    function setUserId(id) {
        userId.value = id;
    }

    return { userId, setUserId };
});
