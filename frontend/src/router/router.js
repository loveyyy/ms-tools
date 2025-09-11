/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2024-09-27 03:52:56
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-11 15:40:56
 * @Description: 
 */
import Main from "@/views/Main.vue";
// 作为Main组件的子页面展示但是不在左侧菜单显示的路由写在otherRouter里
export const otherRouter = {
  path: "/",
  name: "/",
  redirect: "/home",
  component: Main,
  children: [
    {
      path: 'home',
      title: "首页",
      name: 'home',
      component: () => import("@/views/home.vue"),
    },
    {
      path: 'xhs',
      title: "小红书",
      name: 'xhs',
      component: () => import("@/views/xhs/xhs.vue"),
    },
    {
      path: 'xhsDetail',
      title: "小红书详情",
      name: 'xhsDetail',
      component: () => import("@/views/xhs/detail.vue"),
    },
    {
      path: 'xhsUser',
      title: "小红书用户",
      name: 'xhsUser',
      component: () => import("@/views/xhs/user.vue"),
    },
    {
      path: 'dy',
      title: "抖音",
      name: 'dy',
      component: () => import("@/views/dy/dy.vue"),
    },
    {
      path: 'dyDetail',
      title: "抖音详情",
      name: 'dyDetail',
      component: () => import("@/views/dy/detail.vue"),
    },
    // {
    //   path: 'setting',
    //   title: "设置",
    //   name: 'setting',
    //   component: () => import("@/views/dy/dy.vue"),
    // }
  ],
};

// export const page404 = {
//   path: "/:path(.*)*",
//   name: "notFound",
//   meta: {
//     title: "404-页面不存在",
//   },
//   component: () => import("@/views/error-page/404.vue"),
// };
//
// export const page403 = {
//   path: "/403",
//   meta: {
//     title: "403-权限不足",
//   },
//   name: "error-403",
//   component: () => import("@/views/error-page/403.vue"),
// };
//
// export const page500 = {
//   path: "/500",
//   meta: {
//     title: "500-服务端错误",
//   },
//   name: "error-500",
//   component: () => import("@/views/error-page/500.vue"),
// };

// 所有上面定义的路由都要写在下面的routers里
export const routers = [otherRouter];
