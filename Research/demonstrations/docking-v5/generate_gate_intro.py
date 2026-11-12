import math
import subprocess
import sys
import tempfile
from pathlib import Path

import bpy
from mathutils import Vector


ROOT = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT))
import generate_docking_demo as base


VIDEO_PATH = ROOT / "drone-docking-demo-v5-gate-intro.mp4"
BLEND_PATH = ROOT / "drone-docking-demo-v5-gate-intro.blend"

GATE_Y = -14.5
GATE_Z = 2.8
FRAME_END = 144
SENSOR_PLANE_SPACING_M = 0.4
SENSOR_CROSSING_FRAMES = 4
SPEED_ESTIMATE_MPS = SENSOR_PLANE_SPACING_M / (SENSOR_CROSSING_FRAMES / 24)


def torus(name, location, major_radius, minor_radius, mat):
    bpy.ops.mesh.primitive_torus_add(
        major_segments=64,
        minor_segments=12,
        major_radius=major_radius,
        minor_radius=minor_radius,
        location=location,
    )
    obj = bpy.context.object
    obj.name = name
    obj.rotation_euler[0] = math.pi / 2
    obj.data.materials.append(mat)
    return obj


def sphere(name, location, radius, mat):
    bpy.ops.mesh.primitive_uv_sphere_add(
        segments=20, ring_count=12, radius=radius, location=location
    )
    obj = bpy.context.object
    obj.name = name
    obj.data.materials.append(mat)
    return obj


def build_gate(mats):
    for y in (GATE_Y - 0.2, GATE_Y + 0.2):
        torus("Gate structural hoop", (0, y, GATE_Z), 1.95, 0.075, mats["frame"])
        active = torus(
            "CapSense sensing hoop", (0, y, GATE_Z), 1.78, 0.035, mats["capsense_active"]
        )
        base.reveal_after(active, 66 if y < GATE_Y else 70)

        for index in range(8):
            angle = index * math.tau / 8
            x = 1.78 * math.cos(angle)
            z = GATE_Z + 1.78 * math.sin(angle)
            base.cube(
                "Segmented gate electrode",
                (x, y, z),
                (0.24, 0.075, 0.12),
                mats["copper"],
                0.025,
            )
            glow = base.cube(
                "Active gate electrode",
                (x, y, z),
                (0.20, 0.09, 0.15),
                mats["capsense_active"],
                0.025,
            )
            base.reveal_after(glow, 66 if y < GATE_Y else 70)

    confirmed = torus(
        "Passage confirmed halo",
        (0, GATE_Y + 0.24, GATE_Z),
        2.08,
        0.045,
        mats["gate_confirmed"],
    )
    base.reveal_after(confirmed, 92)

    panel = base.cube(
        "Gate status panel",
        (2.9, GATE_Y + 0.15, 2.75),
        (0.95, 0.14, 0.72),
        mats["dark"],
        0.08,
    )
    label = base.add_text(
        "Gate panel label",
        "GATE",
        (2.9, GATE_Y + 0.30, 2.78),
        0.28,
        mats["white"],
        align="CENTER",
    )
    label.data.font = bpy.data.fonts.load(
        "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
    )
    camera_position = Vector((4.0, -3.0, 4.6))
    label.rotation_mode = "QUATERNION"
    label.rotation_quaternion = (
        camera_position - label.location
    ).to_track_quat("Z", "Y")
    sphere("Gate armed indicator", (2.9, GATE_Y + 0.24, 2.57), 0.075, mats["glow"])
    detected = sphere(
        "Gate detected indicator",
        (2.9, GATE_Y + 0.24, 2.57),
        0.105,
        mats["gate_confirmed"],
    )
    base.reveal_after(detected, 92)


def set_drone_path():
    root = bpy.data.objects["Drone motion path"]
    root.animation_data_clear()
    path = (
        (1, (0.0, -21.0, 2.8)),
        (28, (0.0, -19.0, 2.8)),
        (52, (0.0, -16.2, 2.8)),
        (66, (0.0, -14.7, 2.8)),
        (70, (0.0, -14.3, 2.8)),
        (92, (0.0, -12.1, 2.8)),
        (116, (0.0, -9.6, 2.8)),
        (144, (0.0, -7.0, 2.8)),
    )
    for frame, position in path:
        root.location = position
        root.keyframe_insert(data_path="location", frame=frame)
    for frame, angle in ((1, -3.0), (52, -1.0), (92, 0.8), (144, 0.0)):
        root.rotation_euler[2] = math.radians(angle)
        root.keyframe_insert(data_path="rotation_euler", frame=frame)


def compose_video(frame_prefix):
    font = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
    bold_font = "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
    filters = [
        "drawbox=x=40:y=32:w=620:h=76:color=0x07121d@0.88:t=fill",
        f"drawtext=fontfile={bold_font}:text='CAPSENSE PASSAGE GATE':x=64:y=54:fontsize=30:fontcolor=white",
        "drawbox=x=40:y=ih-262:w=550:h=218:color=0x07121d@0.9:t=fill",
        f"drawtext=fontfile={bold_font}:text='SENSOR STATE':x=64:y=main_h-234:fontsize=22:fontcolor=0x9FB3C2",
        f"drawtext=fontfile={font}:text='ARMED':x=64:y=main_h-192:fontsize=30:fontcolor=0x70E8FF:enable='between(t\\,0\\,2.65)'",
        f"drawtext=fontfile={font}:text='ZONES A  >  B':x=64:y=main_h-192:fontsize=30:fontcolor=0xFF9855:enable='between(t\\,2.65\\,3.8)'",
        f"drawtext=fontfile={bold_font}:text='DRONE DETECTED':x=64:y=main_h-192:fontsize=30:fontcolor=0x62F2A5:enable='gte(t,3.8)'",
        f"drawtext=fontfile={font}:text='PASSAGE COUNT   01':x=64:y=main_h-142:fontsize=24:fontcolor=white:enable='gte(t,3.8)'",
        f"drawtext=fontfile={font}:text='SPEED EST.   {SPEED_ESTIMATE_MPS:.1f} m/s':x=64:y=main_h-96:fontsize=24:fontcolor=0x9FEFC2:enable='gte(t,3.8)'",
        f"drawtext=fontfile={font}:text='CONCEPT ONLY  |  SPEED ESTIMATE ILLUSTRATIVE, NOT VALIDATED':x=40:y=main_h-22:fontsize=17:fontcolor=0xC6D4DF",
    ]
    subprocess.run(
        [
            "ffmpeg",
            "-v",
            "error",
            "-framerate",
            "24",
            "-start_number",
            "1",
            "-i",
            f"{frame_prefix}%04d.png",
            "-vf",
            ",".join(filters),
            "-c:v",
            "libx264",
            "-preset",
            "medium",
            "-crf",
            "19",
            "-pix_fmt",
            "yuv420p",
            "-movflags",
            "+faststart",
            "-y",
            str(VIDEO_PATH),
        ],
        check=True,
    )


def main():
    bpy.ops.object.select_all(action="SELECT")
    bpy.ops.object.delete(use_global=False)
    for datablocks in (
        bpy.data.meshes,
        bpy.data.curves,
        bpy.data.materials,
        bpy.data.cameras,
        bpy.data.lights,
    ):
        for datablock in list(datablocks):
            if datablock.users == 0:
                datablocks.remove(datablock)

    mats = {
        "dark": base.material("Gate graphite", (0.035, 0.07, 0.12), 0.65, 0.3),
        "floor": base.material("Flight path floor", (0.025, 0.045, 0.075), 0.25, 0.6),
        "frame": base.material("Gate frame", (0.11, 0.22, 0.31), 0.72, 0.28),
        "drone": base.material("Drone shell", (0.77, 0.83, 0.9), 0.35, 0.3),
        "accent": base.material("Drone blue accent", (0.03, 0.42, 0.72), 0.45, 0.25),
        "prop": base.material("Drone rotors", (0.08, 0.16, 0.22), 0.5, 0.38),
        "rotor_blur": base.translucent_material(
            "Translucent rotor sweep", (0.35, 0.58, 0.72), 0.16
        ),
        "copper": base.material(
            "Gate electrode copper", (0.95, 0.31, 0.08), 0.68, 0.28, ((0.35, 0.06, 0.01), 0.55)
        ),
        "capsense_active": base.material(
            "Active gate sensor", (1.0, 0.28, 0.06), 0.15, 0.24, ((1.0, 0.16, 0.015), 4.0)
        ),
        "glow": base.material(
            "Gate armed light", (0.04, 0.75, 0.85), 0.2, 0.24, ((0.02, 0.75, 1.0), 2.5)
        ),
        "white": base.material("Panel text", (0.85, 0.94, 1.0), 0.0, 0.3),
        "gate_confirmed": base.material(
            "Passage confirmed green",
            (0.08, 0.9, 0.43),
            0.1,
            0.22,
            ((0.02, 1.0, 0.28), 4.0),
        ),
    }

    base.cube("Flight lane floor", (0, -12.0, -0.16), (12, 27, 0.3), mats["floor"], 0.12)
    for x in (-1.9, 1.9):
        base.cube("Flight lane marker", (x, -12.0, 0.005), (0.035, 25, 0.012), mats["frame"])
    build_gate(mats)
    base.build_drone(mats)
    set_drone_path()

    bpy.ops.object.camera_add(location=(4.0, -3.0, 4.6))
    camera = bpy.context.object
    camera.name = "Gate passage camera"
    target = Vector((0, GATE_Y, GATE_Z))
    camera.rotation_euler = (target - camera.location).to_track_quat("-Z", "Y").to_euler()
    camera.data.lens = 29

    for name, location, energy, size, color in (
        ("Gate key light", (2, -11, 10), 1800, 8.0, (0.52, 0.76, 1.0)),
        ("Gate fill light", (-7, -15, 6), 1300, 7.0, (0.16, 0.72, 0.85)),
        ("Gate rim light", (1, -20, 8), 2000, 6.0, (0.15, 0.45, 1.0)),
    ):
        bpy.ops.object.light_add(type="AREA", location=location)
        light = bpy.context.object
        light.name = name
        light.data.energy = energy
        light.data.shape = "DISK"
        light.data.size = size
        light.data.color = color
        light.rotation_euler = (
            Vector((0, GATE_Y, 2.5)) - light.location
        ).to_track_quat("-Z", "Y").to_euler()

    scene = bpy.context.scene
    scene.frame_start = 1
    scene.frame_end = FRAME_END
    scene.render.engine = "BLENDER_EEVEE"
    scene.eevee.taa_render_samples = 8
    scene.render.resolution_x = 1920
    scene.render.resolution_y = 1080
    scene.render.resolution_percentage = 100
    scene.render.fps = 24
    scene.render.image_settings.file_format = "PNG"
    scene.render.image_settings.color_mode = "RGB"
    scene.render.film_transparent = False
    scene.world.color = (0.012, 0.02, 0.04)
    scene.view_settings.view_transform = "AgX"
    scene.camera = camera
    scene.render.filepath = str(ROOT / "drone-docking-demo-v4-gate-frame_")
    scene.frame_set(1)
    bpy.ops.wm.save_as_mainfile(filepath=str(BLEND_PATH))
    with tempfile.TemporaryDirectory(prefix="capsense-gate-frames-") as frame_dir:
        frame_prefix = str(Path(frame_dir) / "frame_")
        scene.render.filepath = frame_prefix
        bpy.ops.render.render(animation=True)
        compose_video(frame_prefix)


if __name__ == "__main__":
    main()
