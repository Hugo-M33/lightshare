import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/models/device.dart';
import '../../../core/providers/devices_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';

class EffectsPanel extends ConsumerWidget {
  final Device device;

  const EffectsPanel({
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
                Icons.auto_awesome,
                color: AppTheme.primaryPurple,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                'Effects',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildEffectButton(
                  context,
                  ref,
                  label: 'Pulse',
                  icon: Icons.favorite,
                  color: AppTheme.primaryPurple,
                  onTap: device.isAvailable
                      ? () => _triggerPulseEffect(ref)
                      : null,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildEffectButton(
                  context,
                  ref,
                  label: 'Breathe',
                  icon: Icons.air,
                  color: AppTheme.accentPink,
                  onTap: device.isAvailable
                      ? () => _triggerBreatheEffect(ref)
                      : null,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildEffectButton(
    BuildContext context,
    WidgetRef ref, {
    required String label,
    required IconData icon,
    required Color color,
    required VoidCallback? onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: color.withValues(alpha: 0.1),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: color.withValues(alpha: 0.3),
          ),
        ),
        child: Column(
          children: [
            Icon(icon, color: color, size: 28),
            const SizedBox(height: 8),
            Text(
              label,
              style: TextStyle(
                color: color,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _triggerPulseEffect(WidgetRef ref) {
    // Use current device color if available, otherwise use a nice purple
    Map<String, dynamic>? color;
    if (device.color != null) {
      color = {
        'hue': device.color!.hue,
        'saturation': device.color!.saturation,
        'kelvin': device.color!.kelvin,
      };
    }

    ref.read(devicesProvider.notifier).pulseEffect(
          device.accountId,
          device.id,
          cycles: 3,
          period: 1.0,
          color: color,
        );
  }

  void _triggerBreatheEffect(WidgetRef ref) {
    // Use current device color if available, otherwise use a nice pink
    Map<String, dynamic>? color;
    if (device.color != null) {
      color = {
        'hue': device.color!.hue,
        'saturation': device.color!.saturation,
        'kelvin': device.color!.kelvin,
      };
    }

    ref.read(devicesProvider.notifier).breatheEffect(
          device.accountId,
          device.id,
          cycles: 3,
          period: 2.0,
          color: color,
        );
  }
}
