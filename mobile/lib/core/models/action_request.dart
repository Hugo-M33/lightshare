/// Action request for executing control actions on devices
class ActionRequest {
  final String action;
  final Map<String, dynamic> parameters;

  ActionRequest({
    required this.action,
    required this.parameters,
  });

  Map<String, dynamic> toJson() {
    return {
      'action': action,
      'parameters': parameters,
    };
  }

  factory ActionRequest.fromJson(Map<String, dynamic> json) {
    return ActionRequest(
      action: json['action'] as String,
      parameters: json['parameters'] as Map<String, dynamic>,
    );
  }

  // Factory constructors for common actions

  /// Create a power action request
  factory ActionRequest.power({
    required bool state,
    double duration = 0.0,
  }) {
    return ActionRequest(
      action: 'power',
      parameters: {
        'state': state ? 'on' : 'off',
        'duration': duration,
      },
    );
  }

  /// Create a brightness action request
  factory ActionRequest.brightness({
    required double level,
    double duration = 0.0,
  }) {
    return ActionRequest(
      action: 'brightness',
      parameters: {
        'level': level,
        'duration': duration,
      },
    );
  }

  /// Create a color action request
  factory ActionRequest.color({
    required double hue,
    required double saturation,
    int? kelvin,
    double duration = 0.0,
  }) {
    final params = {
      'hue': hue,
      'saturation': saturation,
      'duration': duration,
    };
    if (kelvin != null) {
      params['kelvin'] = kelvin.toDouble();
    }
    return ActionRequest(
      action: 'color',
      parameters: params,
    );
  }

  /// Create a temperature action request
  factory ActionRequest.temperature({
    required int kelvin,
    double duration = 0.0,
  }) {
    return ActionRequest(
      action: 'temperature',
      parameters: {
        'kelvin': kelvin.toDouble(),
        'duration': duration,
      },
    );
  }

  /// Create an effect action request
  factory ActionRequest.effect({
    required String name,
    int cycles = 3,
    double period = 1.0,
    Map<String, dynamic>? color,
  }) {
    final params = {
      'name': name,
      'cycles': cycles,
      'period': period,
    };
    if (color != null) {
      params['color'] = color;
    }
    return ActionRequest(
      action: 'effect',
      parameters: params,
    );
  }
}

/// Effect names for device effects
class DeviceEffect {
  static const String pulse = 'pulse';
  static const String breathe = 'breathe';
}
