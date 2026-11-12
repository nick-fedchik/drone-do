import math
import subprocess
import tempfile
from pathlib import Path

import bpy
from mathutils import Vector


ROOT = Path(__file__).resolve().parent
VIDEO_PATH = ROOT / "drone-docking-demo-v4-docking.mp4"
RAW_MAIN_PATH = ROOT / "drone-docking-demo-v4-docking-main.mp4"
RAW_POV_PATH = ROOT / "drone-docking-demo-v4-docking-camera.mp4"
BLEND_PATH = ROOT / "drone-docking-demo-v4-docking.blend"


def material(name, color, metallic=0.0, roughness=0.45, emission=None):
    mat = bpy.data.materials.new(name)
    mat.diffuse_color = (*color, 1.0)
    shader = mat.node_tree.nodes.get("Principled BSDF")
    shader.inputs["Base Color"].default_value = (*color, 1.0)
    shader.inputs["Metallic"].default_value = metallic
    shader.inputs["Roughness"].default_value = roughness
    if emission:
        shader.inputs["Emission Color"].default_value = (*emission[0], 1.0)
        shader.inputs["Emission Strength"].default_value = emission[1]
    return mat


def translucent_material(name, color, alpha):
    mat = material(name, color, 0.1, 0.3)
    mat.diffuse_color = (*color, alpha)
    mat.node_tree.nodes["Principled BSDF"].inputs["Alpha"].default_value = alpha
    mat.surface_render_method = "BLENDED"
    return mat


def cube(name, location, dimensions, mat, bevel=0.0):
    bpy.ops.mesh.primitive_cube_add(size=1, location=location)
    obj = bpy.context.object
    obj.name = name
    obj.dimensions = dimensions
    bpy.ops.object.transform_apply(location=False, rotation=False, scale=True)
    obj.data.materials.append(mat)
    if bevel:
        mod = obj.modifiers.new("Soft edges", "BEVEL")
        mod.width = bevel
        mod.segments = 3
    return obj


def cylinder(name, location, radius, depth, mat, vertices=24):
    bpy.ops.mesh.primitive_cylinder_add(
        vertices=vertices, radius=radius, depth=depth, location=location
    )
    obj = bpy.context.object
    obj.name = name
    obj.data.materials.append(mat)
    return obj


def rod(name, start, end, radius, mat):
    a, b = Vector(start), Vector(end)
    delta = b - a
    obj = cylinder(name, (a + b) / 2, radius, delta.length, mat, vertices=12)
    obj.rotation_mode = "QUATERNION"
    obj.rotation_quaternion = delta.to_track_quat("Z", "Y")
    return obj


def add_text(name, body, location, size, mat, rotation=(0, 0, 0), align="LEFT"):
    data = bpy.data.curves.new(name, "FONT")
    data.body = body
    data.size = size
    data.align_x = align
    data.extrude = 0.001
    obj = bpy.data.objects.new(name, data)
    bpy.context.collection.objects.link(obj)
    obj.location = location
    obj.rotation_euler = rotation
    obj.data.materials.append(mat)
    return obj


def parent_to(obj, parent):
    obj.parent = parent
    obj.matrix_parent_inverse = parent.matrix_world.inverted()


def reveal_after(obj, frame):
    obj.hide_render = True
    obj.keyframe_insert(data_path="hide_render", frame=1)
    obj.keyframe_insert(data_path="hide_render", frame=frame - 1)
    obj.hide_render = False
    obj.keyframe_insert(data_path="hide_render", frame=frame)


def animate_emission(mat, keyframes):
    data_path = 'nodes["Principled BSDF"].inputs["Emission Strength"].default_value'
    strength = mat.node_tree.nodes["Principled BSDF"].inputs["Emission Strength"]
    for frame, value in keyframes:
        strength.default_value = value
        mat.node_tree.keyframe_insert(data_path=data_path, frame=frame)


def build_drone(mats):
    root = bpy.data.objects.new("Drone motion path", None)
    bpy.context.collection.objects.link(root)
    original_objects = set(bpy.context.scene.objects)
    body = cube("Drone body", (0, 0, 0), (0.72, 0.62, 0.24), mats["drone"], 0.09)
    parent_to(body, root)
    top = cube("Battery cover", (0, 0.04, 0.15), (0.38, 0.34, 0.09), mats["accent"], 0.04)
    parent_to(top, root)

    rotor_objects = []
    for x in (-0.72, 0.72):
        for y in (-0.68, 0.68):
            rod("Drone arm", (x * 0.28, y * 0.28, 0), (x, y, 0.02), 0.055, mats["drone"])
            motor = cylinder("Motor", (x, y, 0.06), 0.13, 0.15, mats["dark"])
            parent_to(motor, root)
            hub = cylinder("Rotor hub", (x, y, 0.17), 0.075, 0.07, mats["accent"])
            parent_to(hub, root)
            bpy.ops.mesh.primitive_circle_add(
                vertices=48,
                radius=0.4,
                fill_type="NGON",
                location=(x, y, 0.235),
            )
            blur_disc = bpy.context.object
            blur_disc.name = "Translucent rotor sweep"
            blur_disc.data.materials.append(mats["rotor_blur"])
            parent_to(blur_disc, root)
            blade = cube("Rotor blade", (x, y, 0.21), (0.72, 0.075, 0.025), mats["prop"], 0.025)
            parent_to(blade, root)
            rotor_objects.append(blade)

    for x in (-0.42, 0.42):
        for y in (-0.3, 0.3):
            rod("Landing leg", (x, y, -0.08), (x, y, -0.47), 0.035, mats["drone"])
    for y in (-0.3, 0.3):
        rod("Landing skid", (-0.48, y, -0.47), (0.48, y, -0.47), 0.035, mats["drone"])

    for obj in set(bpy.context.scene.objects) - original_objects:
        if obj.parent is None:
            parent_to(obj, root)

    path = [
        (1, (0.0, -7.0, 2.8)),
        (52, (0.0, -5.2, 2.8)),
        (112, (0.0, -2.0, 2.45)),
        (170, (0.0, 0.2, 2.05)),
        (205, (0.0, 0.2, 1.55)),
        (230, (0.0, 0.2, 1.36)),
        (288, (0.0, 0.2, 1.36)),
    ]
    for frame, position in path:
        root.location = position
        root.keyframe_insert(data_path="location", frame=frame)
    for frame in (1, 288):
        root.rotation_euler[2] = math.radians(4 if frame == 1 else 0)
        root.keyframe_insert(data_path="rotation_euler", frame=frame)

    for blade in rotor_objects:
        for frame in range(1, 289):
            blade.rotation_euler[2] = (frame - 1) * (108 * math.pi / 287)
            blade.keyframe_insert(data_path="rotation_euler", frame=frame)


def build_fan(mats, x, y):
    rotor = bpy.data.objects.new("Vent fan rotor", None)
    bpy.context.collection.objects.link(rotor)
    rotor.location = (x, y, 0.36)
    bpy.context.view_layer.update()

    for angle in range(0, 360, 90):
        radians = math.radians(angle)
        blade = cube(
            "Vent fan blade",
            (x + math.cos(radians) * 0.22, y + math.sin(radians) * 0.22, 0.36),
            (0.72, 0.16, 0.065),
            mats["fan_blade"],
            0.04,
        )
        blade.rotation_euler[2] = radians
        parent_to(blade, rotor)

    bpy.ops.mesh.primitive_torus_add(
        major_segments=48,
        minor_segments=12,
        major_radius=0.64,
        minor_radius=0.05,
        location=(x, y, 0.3),
    )
    active_ring = bpy.context.object
    active_ring.name = "Vent fan active ring"
    active_ring.data.materials.append(mats["fan_glow"])
    for frame, rotation in (
        (1, 0.0),
        (80, 0.0),
        (145, 10 * math.pi),
        (205, 38 * math.pi),
        (225, 46 * math.pi),
        (245, 49 * math.pi),
        (270, 50 * math.pi),
        (288, 50 * math.pi),
    ):
        rotor.rotation_euler[2] = rotation
        rotor.keyframe_insert(data_path="rotation_euler", frame=frame)


def build_hangar(mats):
    floor = cube("Hangar floor", (0, 0, -0.18), (9.5, 9.0, 0.3), mats["floor"], 0.12)
    floor.name = "Hangar floor"
    for x in (-2.9, 2.9):
        cube("Dock side support", (x, 0.25, 0.49), (0.24, 5.3, 0.28), mats["dark"], 0.06)
    for y in (-2.4, 2.9):
        cube("Dock cross support", (0, y, 0.49), (5.8, 0.24, 0.28), mats["dark"], 0.06)
    for x in (-2.7, 2.7):
        for y in (-2.2, 2.7):
            cube("Dock support leg", (x, y, 0.22), (0.25, 0.25, 0.44), mats["frame"], 0.04)

    # Open-front frame: the empty front makes the approach path visible.
    frame_x, front_y, back_y, height = 4.0, -3.5, 4.0, 5.1
    for x in (-frame_x, frame_x):
        for y in (front_y, back_y):
            rod("Structural upright", (x, y, 0), (x, y, height), 0.09, mats["frame"])
        rod("Side roof rail", (x, front_y, height), (x, back_y, height), 0.09, mats["frame"])
        rod("Side lower rail", (x, front_y, 0.08), (x, back_y, 0.08), 0.055, mats["frame"])
    for y in (front_y, back_y):
        rod("Cross roof rail", (-frame_x, y, height), (frame_x, y, height), 0.09, mats["frame"])
    rod("Back mid rail", (-frame_x, back_y, 2.55), (frame_x, back_y, 2.55), 0.045, mats["frame"])
    for x in (-frame_x, frame_x):
        rod("Side mid rail", (x, front_y, 2.55), (x, back_y, 2.55), 0.04, mats["frame"])
    # CapSense electrodes and the glow overlays are conceptual zone indicators.
    for index, (y, activation_frame) in enumerate(
        zip((-2.7, -1.5, -0.3, 0.9), (58, 96, 134, 168)), start=1
    ):
        for x in (-3.82, 3.82):
            electrode = cube(
                f"CapSense electrode Z{index}",
                (x, y, 1.7),
                (0.09, 0.38, 0.16),
                mats["copper"],
                0.025,
            )
            active = cube(
                f"CapSense active glow Z{index}",
                (x, y, 1.7),
                (0.12, 0.32, 0.19),
                mats["capsense_active"],
                0.035,
            )
            reveal_after(active, activation_frame)

    for x in (-2.0, -0.7, 0.7, 2.0):
        electrode = cube("Landing zone electrode", (x, 0.2, 0.755), (0.58, 0.11, 0.018), mats["copper"], 0.025)
        active = cube("Landing electrode active glow", (x, 0.2, 0.765), (0.52, 0.07, 0.02), mats["capsense_active"], 0.025)
        reveal_after(active, 178)

    # Twice the previous bar count, with half the spacing in both directions.
    for index in range(11):
        y = -1.575 + index * 0.35
        cube("Dense landing grate", (0, y, 0.82), (5.35, 0.035, 0.045), mats["grate"], 0.01)
    for index in range(16):
        x = -2.625 + index * 0.35
        cube("Dense landing grate", (x, 0.2, 0.82), (0.035, 3.65, 0.045), mats["grate"], 0.01)

    fan_positions = [(x, y) for x in (-1.35, 1.35) for y in (0.0, 1.35)]
    for x, y in fan_positions:
        cylinder("Fan well", (x, y, 0.13), 0.61, 0.08, mats["dark"], 48)
        bpy.ops.mesh.primitive_torus_add(
            major_segments=48,
            minor_segments=12,
            major_radius=0.59,
            minor_radius=0.07,
            location=(x, y, 0.23),
        )
        housing = bpy.context.object
        housing.name = "Vent fan housing"
        housing.data.materials.append(mats["frame"])
        cylinder("Vent fan hub", (x, y, 0.4), 0.17, 0.14, mats["accent"], 24)
        build_fan(mats, x, y)
    for x in (-3.1, 3.1):
        for y in (-2.7, 3.0):
            light = cylinder("Dock status light", (x, y, 0.24), 0.11, 0.08, mats["glow"], 24)
            light.name = "Dock status light"

    # Forward-facing video and ToF modules mount to the middle rear crossbeam.
    cube("Video camera sensor housing", (0, 3.68, 2.32), (0.48, 0.36, 0.28), mats["dark"], 0.06)
    rod("Video camera mount", (0, 3.95, 2.55), (0, 3.83, 2.44), 0.035, mats["frame"])
    camera_lens = cylinder("Video camera lens", (0, 3.47, 2.32), 0.105, 0.08, mats["glow"], 32)
    camera_lens.rotation_euler[0] = math.pi / 2
    cube("ToF sensor housing", (0.63, 3.68, 2.32), (0.32, 0.32, 0.26), mats["dark"], 0.045)
    rod("ToF sensor mount", (0.63, 3.95, 2.55), (0.63, 3.84, 2.44), 0.03, mats["frame"])
    tof_lens = cylinder("ToF sensor aperture", (0.63, 3.49, 2.32), 0.075, 0.06, mats["accent"], 32)
    tof_lens.rotation_euler[0] = math.pi / 2
    add_text("Camera module label", "CAM", (-0.16, 3.49, 2.45), 0.09, mats["white"], (math.pi / 2, 0, 0))
    add_text("ToF module label", "ToF", (0.51, 3.49, 2.45), 0.09, mats["glow"], (math.pi / 2, 0, 0))

    # Operator interface remains outside the sensor view.
    cube("Operator HMI enclosure", (4.47, -3.32, 2.3), (0.75, 0.16, 1.15), mats["dark"], 0.08)
    cube("HMI display", (4.47, -3.415, 2.3), (0.61, 0.025, 0.91), mats["screen"], 0.035)
    add_text("HMI title", "DOCK", (4.22, -3.439, 2.62), 0.15, mats["glow"], (math.pi / 2, 0, 0))
    add_text("HMI status", "READY", (4.22, -3.439, 2.38), 0.12, mats["white"], (math.pi / 2, 0, 0))
    add_text("HMI charge status", "CHARGING", (4.22, -3.439, 2.12), 0.085, mats["white"], (math.pi / 2, 0, 0))


def make_camera_and_lights(mats):
    bpy.ops.object.camera_add(location=(10.5, -14.8, 14.0))
    camera = bpy.context.object
    camera.name = "Demo camera"
    target = Vector((0, -0.1, 2.1))
    camera.rotation_euler = (target - camera.location).to_track_quat("-Z", "Y").to_euler()
    camera.data.lens = 54
    bpy.context.scene.camera = camera

    bpy.ops.object.camera_add(location=(0, 3.43, 2.32))
    pov_camera = bpy.context.object
    pov_camera.name = "Flight camera POV"
    pov_target = Vector((0, -3.0, 1.8))
    pov_camera.rotation_euler = (pov_target - pov_camera.location).to_track_quat("-Z", "Y").to_euler()
    pov_camera.data.lens = 24

    for name, location, energy, size, color in [
        ("Key light", (1, -6, 10), 1800, 8.0, (0.52, 0.76, 1.0)),
        ("Fill light", (-7, -1, 6), 1300, 7.0, (0.16, 0.72, 0.85)),
        ("Rim light", (1, 6, 8), 2000, 6.0, (0.15, 0.45, 1.0)),
    ]:
        bpy.ops.object.light_add(type="AREA", location=location)
        light = bpy.context.object
        light.name = name
        light.data.energy = energy
        light.data.shape = "DISK"
        light.data.size = size
        light.data.color = color
        light.rotation_euler = (Vector((0, 0, 1.5)) - light.location).to_track_quat("-Z", "Y").to_euler()

    return camera, pov_camera


def compose_video():
    font = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
    bold_font = "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
    filters = [
        "drawbox=x=iw-432:y=0:w=432:h=243:color=0x61d8ff@1:t=5",
        "drawtext=fontfile=" + bold_font + ":text='CAPSENSE DRONE DOCKING':x=48:y=38:fontsize=30:fontcolor=white:box=1:boxcolor=0x07121d@0.88:boxborderw=12",
        "drawbox=x=iw-681:y=ih-141:w=654:h=114:color=0x07121d@0.9:t=fill",
        "drawtext=fontfile=" + font + ":text='Concept only - CapSense is not a drone identifier':x=36:y=main_h-42:fontsize=20:fontcolor=0xC6D4DF",
    ]
    for text, x, y in (
        ("CAPSENSE", "main_w-663", "main_h-122"),
        ("TOF", "main_w-383", "main_h-122"),
        ("VIDEO CAM", "main_w-663", "main_h-72"),
        ("4 FANS", "main_w-383", "main_h-72"),
    ):
        filters.append(
            f"drawtext=fontfile={bold_font}:text='{text}':x={x}:y={y}:"
            "fontsize=23:fontcolor=0x657584"
        )
    active_labels = (
        ("CAPSENSE", "main_w-663", "main_h-122", "0xFF9855", ((2.4, 3.1), (4.0, 4.7), (5.55, 6.25), (7.0, 8.0))),
        ("TOF", "main_w-383", "main_h-122", "0x70E8FF", ((2.0, 2.5), (4.5, 5.0), (7.0, 7.5), (8.5, 9.0))),
        ("VIDEO CAM", "main_w-663", "main_h-72", "0xE7F5FF", ((2.0, 2.5), (4.5, 5.0), (7.0, 7.5), (9.3, 9.8))),
        ("4 FANS", "main_w-383", "main_h-72", "0x70BFFF", ((3.3, 10.8),)),
    )
    for text, x, y, color, intervals in active_labels:
        enable = "+".join(
            f"between(t\\,{start:.2f}\\,{end:.2f})" for start, end in intervals
        )
        pulse = "if(lt(mod(t*4\\,2)\\,1)\\,1\\,0.18)"
        filters.append(
            f"drawtext=fontfile={bold_font}:text='{text}':x={x}:y={y}:"
            f"fontsize=23:fontcolor={color}:alpha='{pulse}':enable='{enable}'"
        )
    filter_graph = (
        "[1:v]scale=432:243[pip];"
        "[0:v][pip]overlay=x=W-w:y=0:shortest=1[composite];"
        "[composite]" + ",".join(filters) + "[out]"
    )
    subprocess.run(
        [
            "ffmpeg",
            "-v",
            "error",
            "-i",
            str(RAW_MAIN_PATH),
            "-i",
            str(RAW_POV_PATH),
            "-filter_complex",
            filter_graph,
            "-map",
            "[out]",
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
            "-shortest",
            "-y",
            str(VIDEO_PATH),
        ],
        check=True,
    )


def encode_sequence(frame_prefix, output_path):
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
            "-c:v",
            "libx264",
            "-preset",
            "medium",
            "-crf",
            "18",
            "-pix_fmt",
            "yuv420p",
            "-movflags",
            "+faststart",
            "-y",
            str(output_path),
        ],
        check=True,
    )


def main():
    bpy.ops.object.select_all(action="SELECT")
    bpy.ops.object.delete(use_global=False)
    for datablocks in (bpy.data.meshes, bpy.data.curves, bpy.data.materials, bpy.data.cameras, bpy.data.lights):
        for datablock in list(datablocks):
            if datablock.users == 0:
                datablocks.remove(datablock)

    mats = {
        "dark": material("Graphite", (0.035, 0.07, 0.12), 0.65, 0.3),
        "floor": material("Hangar floor", (0.025, 0.045, 0.075), 0.35, 0.55),
        "frame": material("Anodized frame", (0.11, 0.22, 0.31), 0.72, 0.28),
        "drone": material("Drone shell", (0.77, 0.83, 0.9), 0.35, 0.3),
        "accent": material("Blue accent", (0.03, 0.42, 0.72), 0.45, 0.25, ((0.01, 0.18, 0.55), 0.4)),
        "prop": material("Rotor and fan blades", (0.08, 0.16, 0.22), 0.5, 0.38),
        "rotor_blur": translucent_material("Translucent rotor sweep", (0.35, 0.58, 0.72), 0.16),
        "fan_blade": material("Visible fan blades", (0.36, 0.65, 0.82), 0.35, 0.25, ((0.04, 0.18, 0.34), 0.8)),
        "grate": material("Landing grate", (0.18, 0.29, 0.37), 0.75, 0.3),
        "copper": material("Illustrative copper electrodes", (0.95, 0.31, 0.08), 0.68, 0.28, ((0.35, 0.06, 0.01), 0.55)),
        "glow": material("Cyan indicators", (0.04, 0.75, 0.85), 0.2, 0.24, ((0.02, 0.75, 1.0), 2.5)),
        "capsense_active": material("Active CapSense zone", (1.0, 0.28, 0.06), 0.15, 0.24, ((1.0, 0.16, 0.015), 4.0)),
        "fan_glow": material("Active fan ring", (0.04, 0.48, 0.95), 0.1, 0.24, ((0.02, 0.3, 1.0), 2.8)),
        "screen": material("HMI screen", (0.005, 0.018, 0.03), 0.15, 0.2),
        "white": material("Overlay white", (0.85, 0.94, 1.0), 0.0, 0.3, ((0.75, 0.9, 1.0), 1.2)),
    }

    build_hangar(mats)
    build_drone(mats)
    main_camera, pov_camera = make_camera_and_lights(mats)

    scene = bpy.context.scene
    scene.frame_start = 1
    scene.frame_end = 288
    scene.render.engine = "BLENDER_EEVEE"
    scene.eevee.taa_render_samples = 6
    scene.render.resolution_x = 1920
    scene.render.resolution_y = 1080
    scene.render.resolution_percentage = 100
    scene.render.fps = 24
    scene.render.image_settings.file_format = "PNG"
    scene.render.image_settings.color_mode = "RGB"
    scene.render.film_transparent = False
    scene.world.color = (0.012, 0.02, 0.04)
    scene.view_settings.view_transform = "AgX"
    scene.render.image_settings.color_mode = "RGB"
    main_camera.data.dof.use_dof = False
    pov_camera.data.dof.use_dof = False

    animate_emission(
        mats["fan_glow"],
        ((1, 0.0), (80, 0.0), (120, 0.7), (170, 2.0), (205, 2.8), (225, 2.6), (245, 1.2), (270, 0.0), (288, 0.0)),
    )
    scene.frame_set(1)
    scene.camera = main_camera
    scene.render.filepath = str(ROOT / "drone-docking-demo-v4-docking-frame_")
    bpy.ops.wm.save_as_mainfile(filepath=str(BLEND_PATH))
    with tempfile.TemporaryDirectory(prefix="capsense-docking-frames-") as frame_dir:
        main_prefix = str(Path(frame_dir) / "main_")
        scene.render.filepath = main_prefix
        bpy.ops.render.render(animation=True)
        encode_sequence(main_prefix, RAW_MAIN_PATH)

        scene.camera = pov_camera
        scene.render.resolution_x = 864
        scene.render.resolution_y = 486
        scene.eevee.taa_render_samples = 8
        pov_prefix = str(Path(frame_dir) / "pov_")
        scene.render.filepath = pov_prefix
        scene.frame_set(1)
        bpy.ops.render.render(animation=True)
        encode_sequence(pov_prefix, RAW_POV_PATH)
        compose_video()


if __name__ == "__main__":
    main()
