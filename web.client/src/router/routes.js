import Home from "../components/pages/Home";
import Room from "../components/pages/Room";
import ColorPicker from "../components/pages/ColorPicker";

export default [
  {
    path: "/",
    name: "home",
    component: Home,
    children: [
      {
        path: "room/:roomId",
        name: "room",
        component: Room,
        children: [
          {
            path: "device/:deviceId/rgb",
            name: "rgb",
            component: ColorPicker
          }
        ]
      }
    ]
  }
  // { path: '/room/:roomId', name: 'room', component: Room },
];
