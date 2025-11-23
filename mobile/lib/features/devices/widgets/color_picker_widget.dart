import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/models/device.dart';
import '../../../core/providers/devices_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';

class ColorPickerWidget extends ConsumerWidget {
  final Device device;

  const ColorPickerWidget({
    super.key,
    required this.device,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return GlassContainer(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(
                Icons.palette,
                color: AppTheme.accentPink,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                'Color',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 20),

          // Hue slider
          _buildHueSlider(ref),
          const SizedBox(height: 20),

          // Saturation slider
          _buildSaturationSlider(ref),
        ],
      ),
    );
  }

  Widget _buildHueSlider(WidgetRef ref) {
    final hue = device.color?.hue ?? 0.0;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Hue',
              style: TextStyle(
                color: AppTheme.textSecondary,
                fontSize: 14,
              ),
            ),
            Text(
              '${hue.toStringAsFixed(0)}°',
              style: const TextStyle(
                color: AppTheme.textPrimary,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            gradient: const LinearGradient(
              colors: [
                Colors.red,
                Colors.yellow,
                Colors.green,
                Colors.cyan,
                Colors.blue,
                Colors.magenta,
                Colors.red,
              ],
            ),
          ),
          child: SliderTheme(
            data: SliderTheme.of(context).copyWith(
              activeTrackColor: Colors.transparent,
              inactiveTrackColor: Colors.transparent,
              thumbColor: Colors.white,
              overlayColor: Colors.white.withValues(alpha: 0.3),
              trackHeight: 32.0,
              thumbShape: const RoundSliderThumbShape(
                enabledThumbRadius: 14.0,
              ),
              overlayShape: const RoundSliderOverlayShape(
                overlayRadius: 24.0,
              ),
            ),
            child: Slider(
              value: hue,
              min: 0.0,
              max: 360.0,
              onChanged: device.isAvailable
                  ? (value) {
                      ref.read(devicesProvider.notifier).setColor(
                            device.accountId,
                            device.id,
                            hue: value,
                            saturation: device.color?.saturation ?? 1.0,
                            kelvin: device.color?.kelvin,
                            duration: 0.5,
                          );
                    }
                  : null,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSaturationSlider(WidgetRef ref) {
    final saturation = device.color?.saturation ?? 1.0;
    final hue = device.color?.hue ?? 0.0;

    // Create gradient from white to the current hue
    final hslColor = HSLColor.fromAHSL(1.0, hue, 1.0, 0.5);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Saturation',
              style: TextStyle(
                color: AppTheme.textSecondary,
                fontSize: 14,
              ),
            ),
            Text(
              '${(saturation * 100).toStringAsFixed(0)}%',
              style: const TextStyle(
                color: AppTheme.textPrimary,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            gradient: LinearGradient(
              colors: [
                Colors.white,
                hslColor.toColor(),
              ],
            ),
          ),
          child: SliderTheme(
            data: SliderTheme.of(context).copyWith(
              activeTrackColor: Colors.transparent,
              inactiveTrackColor: Colors.transparent,
              thumbColor: Colors.white,
              overlayColor: Colors.white.withValues(alpha: 0.3),
              trackHeight: 32.0,
              thumbShape: const RoundSliderThumbShape(
                enabledThumbRadius: 14.0,
              ),
              overlayShape: const RoundSliderOverlayShape(
                overlayRadius: 24.0,
              ),
            ),
            child: Slider(
              value: saturation,
              min: 0.0,
              max: 1.0,
              onChanged: device.isAvailable
                  ? (value) {
                      ref.read(devicesProvider.notifier).setColor(
                            device.accountId,
                            device.id,
                            hue: device.color?.hue ?? 0.0,
                            saturation: value,
                            kelvin: device.color?.kelvin,
                            duration: 0.5,
                          );
                    }
                  : null,
            ),
          ),
        ),
      ],
    );
  }
}
