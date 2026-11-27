import * as THREE from "three";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";


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
  new THREE.BoxGeometry(10000, 2, 10000),
  new THREE.MeshStandardMaterial({ color: 0x3a6b35 })
);
ground.position.y = -1; // górna ściana dokładnie na Y = 0
scene.add(ground);

// ---------- Mechazilla: wieża kratownicowa + ruchome ramiona ----------
// Rakieta celuje w X = 0 (control.LateralDesiredVx), środek masy na ARM_Y.
const TOWER_HEIGHT = 140; // m, całkowita wysokość wieży
const ARM_Y = 125;        // = rocket.TowerHeight w Go (wysokość chwytu)
const TOWER_W = 15;       // m, podstawa 15 × 15
const TOWER_X = 25.5;     // środek wieży; ściana od strony rakiety na X = 18
const FACE_X = TOWER_X - TOWER_W / 2;
const CAR_FRONT = FACE_X - 3.3; // przednia ściana wózka (przed słupami wieży)
const PIVOT_X = CAR_FRONT - 0.7; // oś zawiasu ramienia (słupek przylega do wózka)
const ARM_LEN = 24;       // od zawiasu (X = 14) do 10 m za osią rakiety
const ARM_Z = 6.5;        // zamknięte: odbojnik ~4.8 m od osi, promień kadłuba 4.6
const ARM_OPEN = THREE.MathUtils.degToRad(70);
const ARM_CLOSE_TIME = 3; // s

const towerMat = new THREE.MeshStandardMaterial({ color: 0x4a5a63, metalness: 0.4, roughness: 0.6 });
const armMat = new THREE.MeshStandardMaterial({ color: 0x2b3135, metalness: 0.6, roughness: 0.45 });
const padMat = new THREE.MeshStandardMaterial({ color: 0xb8bcbf, metalness: 0.3, roughness: 0.5 });

// Belka między punktami a i b: jeden wspólny sześcian 1×1×1, rozciągnięty i obrócony.
const unitBox = new THREE.BoxGeometry(1, 1, 1);
const UP = new THREE.Vector3(0, 1, 0);
function addBeam(parent, a, b, t, mat) {
  const d = new THREE.Vector3().subVectors(b, a);
  const m = new THREE.Mesh(unitBox, mat);
  m.scale.set(t, d.length(), t);
  m.position.addVectors(a, b).multiplyScalar(0.5);
  m.quaternion.setFromUnitVectors(UP, d.normalize());
  parent.add(m);
}
const v = (x, y, z) => new THREE.Vector3(x, y, z);

const tower = new THREE.Group();
tower.position.set(TOWER_X, 0, 0);
{
  const h = TOWER_W / 2;
  const corners = [v(-h, 0, -h), v(h, 0, -h), v(h, 0, h), v(-h, 0, h)];
  for (const c of corners) addBeam(tower, c, v(c.x, TOWER_HEIGHT, c.z), 1.6, towerMat);

  const SEG = 10; // m, wysokość jednego segmentu kratownicy
  for (let y0 = 0; y0 < TOWER_HEIGHT; y0 += SEG) {
    const y1 = y0 + SEG;
    for (let i = 0; i < 4; i++) {
      const a = corners[i], b = corners[(i + 1) % 4];
      addBeam(tower, v(a.x, y1, a.z), v(b.x, y1, b.z), 0.8, towerMat); // poziom
      addBeam(tower, v(a.x, y0, a.z), v(b.x, y1, b.z), 0.5, towerMat); // krzyż X
      addBeam(tower, v(b.x, y0, b.z), v(a.x, y1, a.z), 0.5, towerMat);
    }
  }
}
scene.add(tower);

// ---------- plac, fundament wieży, podest startowy ----------
// Ziemia ma górę na Y = 0, więc plac wystaje 0.4 m (bez migotania z ziemią).
const concreteMat = new THREE.MeshStandardMaterial({ color: 0x9c9a94, roughness: 0.95 });
const mountMat = new THREE.MeshStandardMaterial({ color: 0x7d8e98, metalness: 0.3, roughness: 0.6 });

const PAD_H = 0.4;
const pad = new THREE.Mesh(new THREE.BoxGeometry(64, PAD_H, 56), concreteMat);
pad.position.set(11, PAD_H / 2, 0); // X: od −21 do 43; Z: ±28
scene.add(pad);

// Fundament wieży: betonowy blok szerszy o 2.5 m z każdej strony; słupy wchodzą w niego.
const FOUND_H = 3;
const found = new THREE.Mesh(
  new THREE.BoxGeometry(TOWER_W + 5, FOUND_H, TOWER_W + 5), concreteMat
);
found.position.set(TOWER_X, FOUND_H / 2, 0);
scene.add(found);

// Podest startowy (jak orbital launch mount): pierścień z otworem na silniki,
// na 6 pochylonych nogach. Środek na osi rakiety X = 0.
const MOUNT_Y = 18;      // m, spód pierścienia
const MOUNT_R = 8.5;     // m, promień zewnętrzny (kadłub 4.6 m)
const MOUNT_HOLE = 3.5;  // m, promień otworu
{
  // Pierścień = obrót przekroju (r, y) wokół osi Y.
  const profile = [
    new THREE.Vector2(MOUNT_HOLE, 0), new THREE.Vector2(MOUNT_R, 0),
    new THREE.Vector2(MOUNT_R, 2.5), new THREE.Vector2(MOUNT_HOLE, 2.5),
    new THREE.Vector2(MOUNT_HOLE, 0),
  ];
  const ring = new THREE.Mesh(
    new THREE.LatheGeometry(profile, 48),
    new THREE.MeshStandardMaterial({ color: 0x7d8e98, metalness: 0.3, roughness: 0.6, side: THREE.DoubleSide })
  );
  ring.position.set(0, MOUNT_Y, 0);
  scene.add(ring);

  // Nogi: górą przy pierścieniu (r = 7), dołem szerzej (r = 9.5), stoją na placu.
  for (let i = 0; i < 6; i++) {
    const a = (i / 6) * Math.PI * 2 + Math.PI / 6;
    const c = Math.cos(a), s = Math.sin(a);
    addBeam(
      scene,
      v(c * 9.5, PAD_H, s * 9.5),
      v(c * 7, MOUNT_Y + 0.5, s * 7),
      1.8, mountMat
    );
  }
  // Betonowe stopki pod nogami.
  for (let i = 0; i < 6; i++) {
    const a = (i / 6) * Math.PI * 2 + Math.PI / 6;
    const foot = new THREE.Mesh(new THREE.BoxGeometry(3.4, 0.9, 3.4), concreteMat);
    foot.position.set(Math.cos(a) * 9.5, PAD_H + 0.2, Math.sin(a) * 9.5);
    scene.add(foot);
  }
}

// Wózek: pierścień obejmujący wieżę na wysokości ramion (przód + boki + tył).
const CAR_H = 12, CAR_Y = ARM_Y - 2;
{
  const outer = TOWER_W / 2 + 1.3; // poza słupami narożnymi (pół grubości 0.8)
  const front = new THREE.Mesh(new THREE.BoxGeometry(2.5, CAR_H, 2 * outer + 1), armMat);
  front.position.set(CAR_FRONT + 1.25, CAR_Y, 0);
  const back = front.clone();
  back.position.x = TOWER_X + outer + 0.5;
  const sideLen = back.position.x - front.position.x;
  const sideA = new THREE.Mesh(new THREE.BoxGeometry(sideLen, CAR_H, 1), armMat);
  sideA.position.set((front.position.x + back.position.x) / 2, CAR_Y, outer);
  const sideB = sideA.clone();
  sideB.position.z = -outer;
  scene.add(front, back, sideA, sideB);
}

// Ramię wzdłuż −X od osi zawiasu (x = 0). Każdy element dotyka innego:
// słupek zawiasu przy wózku → podpora skośna → kratownica (2 płaszczyzny) → odbojnik.
// inner = −1 lub +1: strona Z zwrócona do rakiety.
function makeArm(inner) {
  const g = new THREE.Group();
  const top = 2, bot = -2, w = 1, step = 3;

  // Słupek zawiasu: od dołu do góry wózka, przylega do jego przedniej ściany.
  addBeam(g, v(0, CAR_Y - ARM_Y - CAR_H / 2, 0), v(0, CAR_Y - ARM_Y + CAR_H / 2, 0), 1.4, armMat);

  for (const z of [-w, w]) {
    addBeam(g, v(0, top, z), v(-ARM_LEN, top, z), 0.7, armMat);
    addBeam(g, v(0, bot, z), v(-ARM_LEN, bot, z), 0.7, armMat);
    let up = true;
    for (let x = 0; x > -ARM_LEN + 0.01; x -= step) {
      const x2 = Math.max(x - step, -ARM_LEN);
      addBeam(g, v(x, bot, z), v(x, top, z), 0.4, armMat);
      addBeam(g, v(x, up ? bot : top, z), v(x2, up ? top : bot, z), 0.35, armMat);
      up = !up;
    }
    addBeam(g, v(-ARM_LEN, bot, z), v(-ARM_LEN, top, z), 0.4, armMat);
    // Podpora od dołu słupka zawiasu do dolnego pasa.
    addBeam(g, v(0, CAR_Y - ARM_Y - CAR_H / 2 + 0.5, z * 0.7), v(-ARM_LEN * 0.4, bot, z), 0.6, armMat);
  }
  // Poprzeczki łączące obie płaszczyzny kratownicy.
  for (let x = 0; x > -ARM_LEN - 0.01; x -= step) {
    addBeam(g, v(x, top, -w), v(x, top, w), 0.3, armMat);
    addBeam(g, v(x, bot, -w), v(x, bot, w), 0.3, armMat);
  }
  // Pokład na górze ramienia.
  const deck = new THREE.Mesh(new THREE.BoxGeometry(ARM_LEN, 0.3, 2 * w + 0.7), padMat);
  deck.position.set(-ARM_LEN / 2, top + 0.5, 0);
  g.add(deck);
  // Odbojnik po stronie rakiety (na odcinku, który obejmuje kadłub).
  const padLen = ARM_LEN * 0.6;
  const pad = new THREE.Mesh(new THREE.BoxGeometry(padLen, 2.6, 0.6), padMat);
  pad.position.set(-ARM_LEN + padLen / 2, 0, inner * (w + 0.6));
  g.add(pad);
  return g;
}

// Dwa ramiona po obu stronach rakiety (±Z). Obrót wokół Y: + odchyla ramię +Z na zewnątrz.
const armA = makeArm(-1);
armA.position.set(PIVOT_X, ARM_Y, ARM_Z);
const armB = makeArm(+1);
armB.position.set(PIVOT_X, ARM_Y, -ARM_Z);
scene.add(armA, armB);

// armK: 0 = otwarte, 1 = zamknięte. Zamyka się przez ARM_CLOSE_TIME po złapaniu.
let armK = 0;
function updateArms(dt, caught) {
  armK += caught ? dt / ARM_CLOSE_TIME : -dt / 1.5;
  armK = THREE.MathUtils.clamp(armK, 0, 1);
  const e = armK * armK * (3 - 2 * armK); // smoothstep: łagodny start i stop
  const a = ARM_OPEN * (1 - e);
  armA.rotation.y = a;
  armB.rotation.y = -a;
}
updateArms(0, false);

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

// flameGroup: przy dyszy, obrót = gimbal.
// Wygląd: krótki jasny rdzeń + pióropusz z miękkich sprite'ów (biało-różowa mgła).
const FLAME_MAX = 90; // m, długość pióropusza przy throttle = 1
const flameGroup = new THREE.Group();
flameGroup.position.y = -25; // po loadzie .glb: dół modelu
flameGroup.visible = false;
rocket.add(flameGroup);

// Miękka okrągła plamka (gradient) jako tekstura cząstki — bez plików graficznych.
function makeSoftTexture() {
  const c = document.createElement("canvas");
  c.width = c.height = 64;
  const g = c.getContext("2d");
  const grad = g.createRadialGradient(32, 32, 0, 32, 32, 32);
  grad.addColorStop(0, "rgba(255,255,255,1)");
  grad.addColorStop(0.4, "rgba(255,255,255,0.5)");
  grad.addColorStop(1, "rgba(255,255,255,0)");
  g.fillStyle = grad;
  g.fillRect(0, 0, 64, 64);
  return new THREE.CanvasTexture(c);
}
const softTex = makeSoftTexture();

// Rdzeń: krótki, bardzo jasny tuż przy dyszy.
const flameCore = new THREE.Mesh(
  new THREE.ConeGeometry(1.2, 1, 20, 1, true),
  new THREE.MeshBasicMaterial({
    color: 0xffffff, transparent: true, opacity: 0.9,
    depthWrite: false, blending: THREE.AdditiveBlending,
  })
);
flameCore.rotation.z = Math.PI; // czubek w dół
flameGroup.add(flameCore);

// Pióropusz: każda cząstka to Sprite (zawsze zwrócony do kamery).
// Rodzi się przy dyszy, leci w dół, rośnie, falują na boki i gaśnie.
const PLUME_N = 160;
const plume = [];
for (let i = 0; i < PLUME_N; i++) {
  const s = new THREE.Sprite(new THREE.SpriteMaterial({
    map: softTex, transparent: true, depthWrite: false,
    blending: THREE.AdditiveBlending, opacity: 0,
  }));
  s.visible = false;
  flameGroup.add(s);
  plume.push({ s, age: 0, life: 0, vy: 0, vx: 0, vz: 0, seed: Math.random() * 100 });
}

const colHot = new THREE.Color(0xffffff);  // przy dyszy
const colCool = new THREE.Color(0xd9b8d6); // dalej: różowo-szara mgła

function spawn(p, th) {
  p.age = 0;
  p.life = 0.5 + Math.random() * 0.4;           // s
  p.vy = -(FLAME_MAX / p.life) * (0.8 + 0.4 * Math.random()) * th;
  p.vx = (Math.random() - 0.5) * 1.5;
  p.vz = (Math.random() - 0.5) * 1.5;
  p.s.position.set((Math.random() - 0.5) * 1, -Math.random() * 2, (Math.random() - 0.5) * 1);
  p.s.visible = true;
}

function killPlume() {
  for (const p of plume) {
    p.age = p.life;
    p.s.visible = false;
    p.s.material.opacity = 0;
  }
}

function updatePlume(dt, th) {
  if (th <= 0.02) {
    killPlume();
    return;
  }
  const t = performance.now() / 1000;
  for (const p of plume) {
    if (p.age >= p.life) {
      if (Math.random() < 0.5) spawn(p, th);
      else { p.s.visible = false; continue; }
    }
    p.age += dt;
    const k = Math.min(p.age / p.life, 1); // 0 = dysza, 1 = koniec
    // Turbulencja: boczne falowanie rośnie z odległością od dyszy.
    const wobble = Math.sin(t * 7 + p.seed) * 3 * k;
    p.s.position.x += (p.vx + wobble) * dt;
    p.s.position.z += p.vz * dt;
    p.s.position.y += p.vy * dt;
    const size = 1.5 + 4.5 * k;             // m: wąsko przy dyszy, szerzej dalej
    p.s.scale.set(size, size, 1);
    p.s.material.color.copy(colHot).lerp(colCool, k);
    p.s.material.opacity = (1 - k) * 0.55 * th;
  }
}

scene.add(rocket);

// modelOffset: korekty skali/orientacji modelu NIE mieszają się z
// rocket.rotation.z = -angle sterowanym z serwera.
// MODEL_Y_OFFSET: przesunięcie modelu wzdłuż osi rakiety [m]; + = wyżej.
// Tylko wizualne — fizyka (środek masy na 125 m przy chwycie) się nie zmienia.
const MODEL_Y_OFFSET = -32;
const modelOffset = new THREE.Group();
modelOffset.position.y = MODEL_Y_OFFSET;
rocket.add(modelOffset);

// Ładowanie .glb jest asynchroniczne: do czasu sukcesu widać body/nose.
new GLTFLoader().load(
  "/models/starship.glb",
  (gltf) => {
    // Box3 liczy w układzie świata — dlatego mierzymy PRZED dodaniem do rakiety,
    // wtedy wynik jest w metrach względem środka modelu.
    const box = new THREE.Box3().setFromObject(gltf.scene);
    const size = box.getSize(new THREE.Vector3());
    flameGroup.position.y = box.min.y + MODEL_Y_OFFSET;
    modelOffset.add(gltf.scene);
    rocket.remove(body, nose);
    body.geometry.dispose();
    nose.geometry.dispose();
    console.log(
      "starship.glb size (m):",
      size.x.toFixed(1), size.y.toFixed(1), size.z.toFixed(1),
      " nozzleY=", flameGroup.position.y.toFixed(1)
    );
  },
  undefined,
  (err) => console.error("starship.glb:", err)
);

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

  if (s.catching && !s.caught) {
    statusEl.textContent = "ŁAPANIE…";
  } else if (s.caught) {
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

let lastT = performance.now();
function animate() {
  requestAnimationFrame(animate);
  const now = performance.now();
  const dt = Math.min((now - lastT) / 1000, 0.05);
  lastT = now;

  if (latest) {
    rocket.position.set(latest.x, latest.y, 0);
    rocket.rotation.z = -latest.angle;

    // Gasimy dopiero po caught/crashed — podczas catching telemetria i silniki dalej żyją.
    const th = (latest.caught || latest.crashed) ? 0 : latest.throttle;
    const len = Math.max(th * FLAME_MAX, 0.001);
    const on = th > 0.02;
    flameGroup.visible = on;
    // Rdzeń: ~15% długości pióropusza, podstawa przy dyszy.
    const coreLen = Math.max(len * 0.15, 0.001);
    flameCore.scale.y = coreLen;
    flameCore.position.y = -0.5 * coreLen;
    flameGroup.rotation.z = -latest.gimbal;
    updatePlume(dt, th);

    follow(rocket.position);
  }
  updateArms(dt, !!(latest?.catching || latest?.caught));

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
