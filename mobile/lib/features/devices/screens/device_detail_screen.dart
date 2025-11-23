import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/models/device.dart';
import '../../../core/providers/devices_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';
import '../widgets/brightness_slider.dart';
import '../widgets/color_picker_widget.dart';
import '../widgets/temperature_slider.dart';
import '../widgets/effects_panel.dart';

class DeviceDetailScreen extends ConsumerWidget {
  final String deviceId;

  const DeviceDetailScreen({
    super.key,
    required this.deviceId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final device = ref.watch(deviceByIdProvider(deviceId));

    if (device == null) {
      return Scaffold(
        body: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                AppTheme.darkBackground,
                AppTheme.primaryPurple.withValues(alpha: 0.2),
                AppTheme.accentPink.withValues(alpha: 0.1),
                AppTheme.darkBackground,
              ],
            ),
          ),
          child: const Center(
            child: CircularProgressIndicator(
              color: AppTheme.primaryPurple,
            ),
          ),
        ),
      );
    }

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              AppTheme.darkBackground,
              AppTheme.primaryPurple.withValues(alpha: 0.2),
              AppTheme.accentPink.withValues(alpha: 0.1),
              AppTheme.darkBackground,
            ],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              // App bar
              Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    IconButton(
                      onPressed: () => context.pop(),
                      icon: const Icon(Icons.arrow_back),
                      style: IconButton.styleFrom(
                        backgroundColor:
                            AppTheme.cardBackground.withValues(alpha: 0.5),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            device.displayName,
                            style: Theme.of(context).textTheme.headlineSmall,
                            overflow: TextOverflow.ellipsis,
                          ),
                          Text(
                            device.providerDisplayName,
                            style: Theme.of(context).textTheme.bodySmall,
                          ),
                        ],
                      ),
                    ),
                    // Status badge
                    _buildStatusBadge(device),
                  ],
                ),
              ),

              // Content
              Expanded(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      // Device icon and power toggle
                      _buildDeviceHeader(context, ref, device),
                      const SizedBox(height: 24),

                      // Brightness control
                      BrightnessSlider(device: device),
                      const SizedBox(height: 16),

                      // Color controls
                      if (device.supportsColor) ...[
                        ColorPickerWidget(device: device),
                        const SizedBox(height: 16),
                      ],

                      // Temperature control
                      if (device.supportsTemperature) ...[
                        TemperatureSlider(device: device),
                        const SizedBox(height: 16),
                      ],

                      // Effects
                      if (device.supportsEffects) ...[
                        EffectsPanel(device: device),
                        const SizedBox(height: 16),
                      ],

                      // Device info
                      _buildDeviceInfo(context, device),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildDeviceHeader(
    BuildContext context,
    WidgetRef ref,
    Device device,
  ) {
    return GlassContainer(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          // Device icon
          Container(
            padding: const EdgeInsets.all(32),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [
                  device.isOn
                      ? AppTheme.primaryPurple.withValues(alpha: 0.3)
                      : AppTheme.cardBackground.withValues(alpha: 0.3),
                  device.isOn
                      ? AppTheme.accentPink.withValues(alpha: 0.3)
                      : AppTheme.cardBackground.withValues(alpha: 0.3),
                ],
              ),
              shape: BoxShape.circle,
              boxShadow: device.isOn
                  ? [
                      BoxShadow(
                        color: AppTheme.primaryPurple.withValues(alpha: 0.3),
                        blurRadius: 40,
                        spreadRadius: 10,
                      ),
                    ]
                  : null,
            ),
            child: Icon(
              Icons.lightbulb,
              size: 64,
              color: device.isOn
                  ? AppTheme.primaryPurple
                  : AppTheme.textSecondary,
            ),
          ),
          const SizedBox(height: 24),

          // Power toggle button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: device.isAvailable
                  ? () {
                      ref.read(devicesProvider.notifier).setPower(
                            device.accountId,
                            device.id,
                            state: !device.isOn,
                          );
                    }
                  : null,
              icon: Icon(device.isOn ? Icons.power : Icons.power_off),
              label: Text(device.isOn ? 'Turn Off' : 'Turn On'),
              style: ElevatedButton.styleFrom(
                backgroundColor: device.isOn
                    ? AppTheme.primaryPurple
                    : AppTheme.cardBackground,
                foregroundColor: AppTheme.textPrimary,
                padding: const EdgeInsets.symmetric(vertical: 20),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                elevation: 0,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDeviceInfo(BuildContext context, Device device) {
    return GlassContainer(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(
                Icons.info_outline,
                color: AppTheme.primaryPurple,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                'Device Information',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          _buildInfoRow('ID', device.id),
          if (device.group != null) _buildInfoRow('Group', device.group!.name),
          if (device.location != null)
            _buildInfoRow('Location', device.location!.name),
          _buildInfoRow('Connected', device.connected ? 'Yes' : 'No'),
          _buildInfoRow('Reachable', device.reachable ? 'Yes' : 'No'),
          _buildInfoRow(
            'Capabilities',
            device.capabilities.join(', '),
          ),
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(
              color: AppTheme.textSecondary,
              fontSize: 14,
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              color: AppTheme.textPrimary,
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusBadge(Device device) {
    final isAvailable = device.isAvailable;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: isAvailable
            ? Colors.green.withValues(alpha: 0.2)
            : Colors.red.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isAvailable
              ? Colors.green.withValues(alpha: 0.5)
              : Colors.red.withValues(alpha: 0.5),
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 8,
            height: 8,
            decoration: BoxDecoration(
              color: isAvailable ? Colors.green : Colors.red,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            isAvailable ? 'Online' : 'Offline',
            style: TextStyle(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              color: isAvailable ? Colors.green : Colors.red,
            ),
          ),
        ],
      ),
    );
  }
}
