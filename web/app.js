import * as THREE from "three";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";

const log = document.getElementById("log");
const statusEl = document.getElementById("status");

// ---------- scena, kamera, renderer ----------
const scene = new THREE.Scene();
scene.background = new THREE.Color(0x87a8d0);

const camera = new THREE.PerspectiveCamera(
  50, window.innerWidth / window.innerHeight, 1, 5000
);
camera.position.set(0, 100, 400);

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

// Obracanie/zoom myszką wokół punktu, na który patrzy kamera (target).
const controls = new OrbitControls(camera, renderer.domElement);
controls.target.set(0, 100, 0);
controls.minDistance = 50;
controls.maxDistance = 2000;
controls.update();

// ---------- światła (materiały Standard bez światła byłyby czarne) ----------
scene.add(new THREE.AmbientLight(0xffffff, 0.6));
const sun = new THREE.DirectionalLight(0xffffff, 1.0);
sun.position.set(200, 500, 300);
scene.add(sun);

// ---------- ziemia i wieża ----------
const ground = new THREE.Mesh(
  new THREE.BoxGeometry(2000, 2, 200),
  new THREE.MeshStandardMaterial({ color: 0x3a6b35 })
);
ground.position.y = -1; // górna ściana dokładnie na Y = 0
scene.add(ground);

const TOWER_HEIGHT = 90; // = rocket.TowerHeight w Go (wysokość ramion)
const towerMat = new THREE.MeshStandardMaterial({ color: 0x777777 });

const column = new THREE.Mesh(new THREE.BoxGeometry(6, 100, 6), towerMat);
column.position.set(12, 50, 0); // obok rakiety (która ma być w X = 0)
scene.add(column);

const arm = new THREE.Mesh(new THREE.BoxGeometry(4.5, 2, 4), towerMat);
arm.position.set(6.75, TOWER_HEIGHT, 0); // ramię na wysokości chwytu
scene.add(arm);

// ---------- rakieta: kadłub, dziób, płomień ----------
const rocket = new THREE.Group();

const body = new THREE.Mesh(
  new THREE.CylinderGeometry(4.5, 4.5, 50, 24),
  new THREE.MeshStandardMaterial({ color: 0xcccccc })
);
rocket.add(body);

const nose = new THREE.Mesh(
  new THREE.ConeGeometry(4.5, 10, 24),
  new THREE.MeshStandardMaterial({ color: 0xdd3333 })
);
nose.position.y = 30;
rocket.add(nose);

// flameGroup siedzi przy dyszy i obraca się razem z gimbalem;
// flame to stożek skalowany długością w dół od dyszy.
const FLAME_MAX = 30; // m, długość płomienia przy throttle = 1
const flameGroup = new THREE.Group();
flameGroup.position.y = -25;
const flame = new THREE.Mesh(
  new THREE.ConeGeometry(3.5, 1, 16),
  new THREE.MeshBasicMaterial({ color: 0xffaa33 })
);
flame.rotation.z = Math.PI; // czubek w dół
flameGroup.add(flame);
rocket.add(flameGroup);

scene.add(rocket);

// ---------- stan z serwera ----------
let latest = null;

function updateHud(s) {
  log.textContent =
    `t=${s.t.toFixed(1)} s\n` +
    `x=${s.x.toFixed(1)}  y=${s.y.toFixed(1)} m\n` +
    `vx=${s.vx.toFixed(1)}  vy=${s.vy.toFixed(1)} m/s\n` +
    `kąt=${s.angle.toFixed(3)}  gimbal=${s.gimbal.toFixed(3)} rad\n` +
    `ciąg=${(s.throttle * 100).toFixed(0)}%  paliwo=${s.fuel.toFixed(0)} kg\n` +
    `aero=${s.aero.toFixed(0)} N`;

  if (s.caught) {
    statusEl.textContent = "ZŁAPANA";
  } else if (s.crashed) {
    statusEl.textContent = "ROZBITA";
  } else {
    statusEl.textContent = "";
  }
}

// Kamera podąża za rakietą: przesuwamy target i kamerę o tę samą różnicę,
// więc OrbitControls zachowuje Twój kąt i odległość.
function follow(pos) {
  const d = pos.clone().sub(controls.target);
  camera.position.add(d);
  controls.target.copy(pos);
}

function animate() {
  requestAnimationFrame(animate);

  if (latest) {
    rocket.position.set(latest.x, latest.y, 0);
    rocket.rotation.z = -latest.angle;

    const len = latest.throttle * FLAME_MAX;
    flame.visible = latest.throttle > 0.02;
    flame.scale.y = Math.max(len, 0.001);
    flame.position.y = -0.5 * len; // podstawa stożka zawsze przy dyszy
    flameGroup.rotation.z = -latest.gimbal;

    follow(rocket.position);
  }

  controls.update();
  renderer.render(scene, camera);
}
animate();

window.addEventListener("resize", () => {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
});

// ---------- WebSocket ----------
// ws:// lub wss:// zależnie od tego, czy strona jest po http czy https
const proto = location.protocol === "https:" ? "wss" : "ws";
const ws = new WebSocket(`${proto}://${location.host}/ws`);

ws.onopen = () => {
  log.textContent = "połączono";
};

ws.onmessage = (e) => {
  latest = JSON.parse(e.data);
  updateHud(latest);
};

ws.onclose = () => {
  log.textContent += "\nrozłączono";
};

// ---------- przyciski ----------
function send(type) {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type }));
  }
}
document.getElementById("btn-start").onclick = () => send("start");
document.getElementById("btn-pause").onclick = () => send("pause");
document.getElementById("btn-reset").onclick = () => send("reset");
