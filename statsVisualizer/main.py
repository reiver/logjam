import json
import re
import plotly.graph_objects as go
from plotly.subplots import make_subplots

# Open the text file
with open('statsVisualizer\data.txt') as file:
    # Read the contents of the file
    file_contents = file.read()

    # Use regular expressions to find JSON elements
    json_elements = re.findall(r'\{.*?\}', file_contents, re.DOTALL)

    # Initialize lists to store data for visualization
    audio_jitter = []
    audio_packets_lost = []
    cum_available_outgoing_bitrate = []
    video_jitter = []
    video_packets_lost = []
    video_frames_per_second = []
    outbound_fps = []
    inbound_fps = []
    remote_inbound_fps = []
    source_fps = []

    # Extract data from JSON elements
    for element in json_elements:
        try:
            json_data = json.loads(element)

            type = json_data.get('type')
            kind = json_data.get('kind')

            jitter = json_data.get('jitter')
            packets_lost = json_data.get('packetsLost')
            available_outgoing_bitrate = json_data.get('availableOutgoingBitrate')
            frames_per_second = json_data.get('framesPerSecond')

            cum_available_outgoing_bitrate.append(available_outgoing_bitrate)

            if kind == 'audio':
                audio_jitter.append(jitter)
                audio_packets_lost.append(packets_lost)
            elif kind == 'video':
                video_jitter.append(jitter)
                video_packets_lost.append(packets_lost)

                if type == 'outbound-rtp':
                    outbound_fps.append(frames_per_second)
                elif type == 'media-source':
                    source_fps.append(frames_per_second)
                elif type == 'inbound-rtp':
                    inbound_fps.append(frames_per_second)
                elif type == 'remote-inbound-rtp':
                    remote_inbound_fps.append(frames_per_second)

        except json.JSONDecodeError as e:
            print(f"Error parsing JSON: {e}")

    # Create subplots for the four graphs
    fig = make_subplots(rows=1, cols=4, subplot_titles=("Jitter", "Packets Lost", "Available Outgoing Bitrate", "FPS"))

    # Add traces to the subplots
    fig.add_trace(go.Box(y=audio_jitter, name='Audio'), row=1, col=1)
    fig.add_trace(go.Box(y=audio_packets_lost, name='Audio'), row=1, col=2)
    fig.add_trace(go.Box(y=cum_available_outgoing_bitrate, name='Bitrate'), row=1, col=3)
    fig.add_trace(go.Box(y=video_jitter, name='Video'), row=1, col=1)
    fig.add_trace(go.Box(y=video_packets_lost, name='Video'), row=1, col=2)
    fig.add_trace(go.Box(y=source_fps, name='Source'), row=1, col=4)
    fig.add_trace(go.Box(y=outbound_fps, name='Outbound'), row=1, col=4)
    fig.add_trace(go.Box(y=inbound_fps, name='Inbound'), row=1, col=4)
    fig.add_trace(go.Box(y=remote_inbound_fps, name='Remote Inbound'), row=1, col=4)



    # Update layout and display the figure
    fig.update_layout(showlegend=False)
    fig.show()
