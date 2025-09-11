import { routers } from "./router";
import { createRouter, createWebHashHistory } from "vue-router";

export const router = new createRouter({
  history: createWebHashHistory(),
  routes: routers,
});

router.beforeEach((to, from, next) => {
  console.log(to,from,next)
  next();
  // const name = to.name;
  // if (!Cookies.get("userInfoManager") && name !== "login") {
  //   // 判断是否已经登录且前往的页面不是登录页
  //   next({
  //     name: "login",
  //   });
  // } else if (Cookies.get("userInfoManager") && name === "login") {
  //   // 判断是否已经登录且前往的是登录页
  //   Util.title();
  //   next({
  //     name: "home_index",
  //   });
  // } else {
  //   Util.toDefaultPage([...routers], name, router, next);
  // }
});

router.afterEach((to) => {
  window.scrollTo(0, 0);
});
