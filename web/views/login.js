import { api, auth } from '/app.js';

const { ref } = Vue;
const { useRouter } = VueRouter;

export const Login = {
  template: `
    <article>
      <header><strong>Login</strong></header>
      <form @submit.prevent="submit">
        <label>Email
          <input v-model="email" type="email" autocomplete="email" required />
        </label>
        <label>Password
          <input v-model="password" type="password" autocomplete="current-password" required />
        </label>
        <button :aria-busy="loading" type="submit">Sign in</button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
      <footer><router-link to="/register">Need an account?</router-link></footer>
    </article>
  `,
  setup() {
    const router = useRouter();
    const email = ref('');
    const password = ref('');
    const error = ref('');
    const loading = ref(false);

    const submit = async () => {
      loading.value = true;
      error.value = '';
      try {
        const out = await api('/api/login', { method: 'POST', body: { email: email.value, password: password.value } });
        auth.setToken(out.sessionToken);
        router.push('/me');
      } catch (e) {
        error.value = e.message;
      } finally {
        loading.value = false;
      }
    };
    return { email, password, error, loading, submit };
  },
};
