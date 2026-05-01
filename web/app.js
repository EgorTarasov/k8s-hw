import { Login } from '/views/login.js';
import { Register } from '/views/register.js';
import { Me } from '/views/me.js';
import { Orders } from '/views/orders.js';
import { Order } from '/views/order.js';

const { createApp, reactive } = Vue;
const { createRouter, createWebHashHistory } = VueRouter;

export const API_BASE = window.SHOPX_API_BASE ?? '';

export const auth = reactive({
  token: localStorage.getItem('session_token') || '',
  setToken(t) {
    this.token = t || '';
    if (t) localStorage.setItem('session_token', t);
    else localStorage.removeItem('session_token');
  },
});

export async function api(path, { method = 'GET', body } = {}) {
  const headers = { 'Content-Type': 'application/json' };
  if (auth.token) headers['X-Session-Token'] = auth.token;

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  let data = null;
  const text = await res.text();
  if (text) {
    try { data = JSON.parse(text); } catch { data = { message: text }; }
  }
  if (!res.ok) {
    const err = new Error((data && data.message) || `HTTP ${res.status}`);
    err.status = res.status;
    err.code = data && data.code;
    throw err;
  }
  return data;
}

const routes = [
  { path: '/', redirect: () => (auth.token ? '/me' : '/login') },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/me', component: Me, meta: { requiresAuth: true } },
  { path: '/orders', component: Orders, meta: { requiresAuth: true } },
  { path: '/orders/:id', component: Order, meta: { requiresAuth: true }, props: true },
];

const router = createRouter({ history: createWebHashHistory(), routes });

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !auth.token) return { path: '/login' };
});

const App = {
  setup() {
    const logout = () => {
      auth.setToken('');
      router.push('/login');
    };
    return { auth, logout };
  },
};

createApp(App).use(router).mount('#app');
