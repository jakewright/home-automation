import Home from "../pages/home";
import Room from "../pages/room";
import NotFound from "../pages/404";
import ColorPicker from "../components/pages/ColorPicker";

export default [
  {
    path: "/",
    name: "home",
    component: Home
  },
  {
    path: "/room/:roomId",
    name: "room",
    component: Room,
    children: [
      {
        path: "device/:deviceId/rgb",
        name: "rgb",
        component: ColorPicker
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: "not-found",
    component: NotFound,
  }
];
