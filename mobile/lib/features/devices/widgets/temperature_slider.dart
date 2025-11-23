import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/models/device.dart';
import '../../../core/providers/devices_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';

class TemperatureSlider extends ConsumerWidget {
  final Device device;

  const TemperatureSlider({
    super.key,
    required this.device,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final kelvin = device.color?.kelvin ?? 3500;

    return GlassContainer(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  const Icon(
                    Icons.thermostat,
                    color: AppTheme.accentPink,
                    size: 20,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'Temperature',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                ],
              ),
              Text(
                '${kelvin}K',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      color: AppTheme.accentPink,
                      fontWeight: FontWeight.bold,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            _getTemperatureName(kelvin),
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: AppTheme.textSecondary,
                ),
          ),
          const SizedBox(height: 16),
          Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              gradient: const LinearGradient(
                colors: [
                  Color(0xFFFFB46B), // Warm orange (1500K)
                  Color(0xFFFFD6AA), // Warm white (2700K)
                  Color(0xFFFFF4E6), // Neutral white (4000K)
                  Color(0xFFCCE6FF), // Cool white (5000K)
                  Color(0xFFB3D9FF), // Daylight (6500K)
                  Color(0xFF99CCFF), // Cool daylight (9000K)
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
                value: kelvin.toDouble(),
                min: 1500.0,
                max: 9000.0,
                divisions: 75, // 100K steps
                onChanged: device.isAvailable
                    ? (value) {
                        ref.read(devicesProvider.notifier).setTemperature(
                              device.accountId,
                              device.id,
                              kelvin: value.toInt(),
                              duration: 0.5,
                            );
                      }
                    : null,
              ),
            ),
          ),
          const SizedBox(height: 8),
          const Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Warm',
                style: TextStyle(
                  color: AppTheme.textSecondary,
                  fontSize: 12,
                ),
              ),
              Text(
                'Cool',
                style: TextStyle(
                  color: AppTheme.textSecondary,
                  fontSize: 12,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _getTemperatureName(int kelvin) {
    if (kelvin < 2000) {
      return 'Candle';
    } else if (kelvin < 2700) {
      return 'Warm White';
    } else if (kelvin < 3000) {
      return 'Soft White';
    } else if (kelvin < 4000) {
      return 'Neutral White';
    } else if (kelvin < 5000) {
      return 'Cool White';
    } else if (kelvin < 6500) {
      return 'Daylight';
    } else {
      return 'Cool Daylight';
    }
  }
}
