import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/models/device.dart';
import '../../../core/providers/devices_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';

class BrightnessSlider extends ConsumerWidget {
  final Device device;

  const BrightnessSlider({
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
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  const Icon(
                    Icons.brightness_6,
                    color: AppTheme.primaryPurple,
                    size: 20,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'Brightness',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                ],
              ),
              Text(
                '${(device.brightness * 100).toStringAsFixed(0)}%',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      color: AppTheme.primaryPurple,
                      fontWeight: FontWeight.bold,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          SliderTheme(
            data: SliderTheme.of(context).copyWith(
              activeTrackColor: AppTheme.primaryPurple,
              inactiveTrackColor:
                  AppTheme.primaryPurple.withValues(alpha: 0.2),
              thumbColor: AppTheme.primaryPurple,
              overlayColor: AppTheme.primaryPurple.withValues(alpha: 0.3),
              trackHeight: 6.0,
              thumbShape: const RoundSliderThumbShape(
                enabledThumbRadius: 12.0,
              ),
              overlayShape: const RoundSliderOverlayShape(
                overlayRadius: 24.0,
              ),
            ),
            child: Slider(
              value: device.brightness,
              min: 0.0,
              max: 1.0,
              onChanged: device.isAvailable
                  ? (value) {
                      // Debouncing is handled in the provider
                      ref.read(devicesProvider.notifier).setBrightness(
                            device.accountId,
                            device.id,
                            level: value,
                            duration: 0.5,
                          );
                    }
                  : null,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Icon(
                Icons.brightness_low,
                color: AppTheme.textSecondary,
                size: 16,
              ),
              const Icon(
                Icons.brightness_high,
                color: AppTheme.textSecondary,
                size: 16,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
